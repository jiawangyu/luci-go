// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package archive implements a storage.Storage instance that retrieves logs
// from a Google Storage archive.
//
// This is a special implementation of storage.Storage, and does not fully
// conform to the API expecations. Namely:
//	- It is read-only. Mutation methods will return storage.ErrReadOnly.
//	- Storage methods ignore the supplied Path argument, instead opting for
//	  the archive configured in its Options.
package archive

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"sync"

	"github.com/luci/luci-go/common/config"
	"github.com/luci/luci-go/common/data/recordio"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/common/gcloud/gs"
	"github.com/luci/luci-go/common/iotools"
	log "github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/logdog/api/logpb"
	"github.com/luci/luci-go/logdog/common/storage"
	"github.com/luci/luci-go/logdog/common/types"

	cloudStorage "cloud.google.com/go/storage"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
)

const (
	// maxStreamRecordSize is the maximum record size we're willing to read from
	// our archived log stream. This will help prevent out-of-memory errors if the
	// arhived log stream is malicious or corrupt.
	//
	// 16MB is larger than the maximum log entry size
	maxStreamRecordSize = 2 * types.MaxLogEntryDataSize
)

// Options is the set of configuration options for this Storage instance.
//
// Unlike other Storage instances, this is bound to a single archived stream.
// Project and Path parameters in requests will be ignored in favor of the
// Google Storage URLs.
type Options struct {
	// IndexURL is the Google Storage URL for the stream's index.
	IndexURL string
	// StreamURL is the Google Storage URL for the stream's entries.
	StreamURL string

	// Client is the HTTP client to use for authentication.
	//
	// Closing this Storage instance does not close the underlying Client.
	Client gs.Client
}

type storageImpl struct {
	*Options
	context.Context

	streamPath gs.Path
	indexPath  gs.Path

	indexMu sync.Mutex
	index   *logpb.LogIndex
}

// New instantiates a new Storage instance, bound to the supplied Options.
func New(ctx context.Context, o Options) (storage.Storage, error) {
	s := storageImpl{
		Options: &o,
		Context: ctx,

		streamPath: gs.Path(o.StreamURL),
		indexPath:  gs.Path(o.IndexURL),
	}

	if !s.streamPath.IsFullPath() {
		return nil, fmt.Errorf("invalid stream URL: %q", s.streamPath)
	}
	if s.indexPath != "" && !s.indexPath.IsFullPath() {
		return nil, fmt.Errorf("invalid index URL: %v", s.indexPath)
	}

	return &s, nil
}

func (s *storageImpl) Close() {}

func (s *storageImpl) Config(storage.Config) error  { return storage.ErrReadOnly }
func (s *storageImpl) Put(storage.PutRequest) error { return storage.ErrReadOnly }

func (s *storageImpl) Get(req storage.GetRequest, cb storage.GetCallback) error {
	idx, err := s.getIndex()
	if err != nil {
		return err
	}

	// Identify the byte offsets that we want to fetch from the entries stream.
	st := buildGetStrategy(&req, idx)
	if st == nil {
		// No more records to read.
		return nil
	}

	switch err := s.getLogEntriesIter(st, cb); errors.Unwrap(err) {
	case nil, io.EOF:
		// We hit the end of our log stream.
		return nil

	case cloudStorage.ErrObjectNotExist, cloudStorage.ErrBucketNotExist:
		return storage.ErrDoesNotExist

	default:
		return errors.Annotate(err).Reason("failed to read log stream").Err()
	}
}

// getLogEntriesImpl retrieves log entries from archive until complete.
func (s *storageImpl) getLogEntriesIter(st *getStrategy, cb storage.GetCallback) error {
	// Get our maximum byte limit. If we are externally constrained via MaxBytes,
	// apply that limit too.
	// Get an archive reader.
	var (
		offset = st.startOffset
		length = st.length()
	)

	log.Fields{
		"offset": offset,
		"length": length,
		"path":   s.streamPath,
	}.Debugf(s, "Creating stream reader for range.")
	storageReader, err := s.Client.NewReader(s.streamPath, int64(offset), length)
	if err != nil {
		log.WithError(err).Errorf(s, "Failed to create stream Reader.")
		return errors.Annotate(err).Reason("failed to create stream Reader").Err()
	}
	defer func() {
		if tmpErr := storageReader.Close(); tmpErr != nil {
			// (Non-fatal)
			log.WithError(tmpErr).Warningf(s, "Error closing stream Reader.")
		}
	}()

	// Count how many bytes we've read.
	cr := iotools.CountingReader{Reader: storageReader}

	// Iteratively update our strategy's start offset each time we read a complete
	// frame.
	var (
		rio       = recordio.NewReader(&cr, maxStreamRecordSize)
		buf       bytes.Buffer
		remaining = st.count
	)
	for {
		// Reset the count so we know how much we read for this frame.
		cr.Count = 0

		sz, r, err := rio.ReadFrame()
		if err != nil {
			return errors.Annotate(err).Reason("failed to read frame").Err()
		}

		buf.Reset()
		buf.Grow(int(sz))

		switch amt, err := buf.ReadFrom(r); {
		case err != nil:
			log.Fields{
				log.ErrorKey:  err,
				"frameOffset": offset,
				"frameSize":   sz,
			}.Errorf(s, "Failed to read frame data.")
			return errors.Annotate(err).Reason("failed to read frame data").Err()

		case amt != sz:
			// If we didn't buffer the complete frame, we hit a premature EOF.
			return errors.Annotate(io.EOF).Reason("incomplete frame read").Err()
		}

		// If we read from offset 0, the first frame will be the log stream's
		// descriptor, which we can discard.
		discardFrame := (offset == 0)
		offset += uint64(cr.Count)
		if discardFrame {
			continue
		}

		// Punt this log entry to our callback, if appropriate.
		entry := storage.MakeEntry(buf.Bytes(), -1)
		switch idx, err := entry.GetStreamIndex(); {
		case err != nil:
			log.Fields{
				log.ErrorKey:  err,
				"frameOffset": offset,
				"frameSize":   sz,
			}.Errorf(s, "Failed to get log entry index.")
			return errors.Annotate(err).Reason("failed to get log entry index").Err()

		case idx < st.startIndex:
			// Skip this entry, as it's before the first requested entry.
			continue
		}

		// We want to punt this entry, but we also want to re-use our Buffer. Clone
		// its data so it is independent.
		entry.D = make([]byte, len(entry.D))
		copy(entry.D, buf.Bytes())
		if !cb(entry) {
			return nil
		}

		// Enforce our limit, if one is supplied.
		if remaining > 0 {
			remaining--
			if remaining == 0 {
				return nil
			}
		}
	}
}

func (s *storageImpl) Tail(project config.ProjectName, path types.StreamPath) (*storage.Entry, error) {
	idx, err := s.getIndex()
	if err != nil {
		return nil, err
	}

	// Get the offset that is as close to our tail record as possible. If we know
	// what that index is (from "idx"), we can request it directly. Otherwise, we
	// will get as close as possible and read forwards from there.
	req := storage.GetRequest{}
	switch {
	case idx.LastStreamIndex > 0:
		req.Index = types.MessageIndex(idx.LastStreamIndex)
		req.Limit = 1

	case len(idx.Entries) > 0:
		req.Index = types.MessageIndex(idx.Entries[len(idx.Entries)-1].StreamIndex)
	}

	// Build a Get strategy for our closest-to-Tail index.
	st := buildGetStrategy(&req, idx)
	if st == nil {
		return nil, storage.ErrDoesNotExist
	}

	// Read forwards to EOF. Retain the last entry that we read.
	var lastEntry *storage.Entry
	err = s.Get(req, func(e *storage.Entry) bool {
		lastEntry = e

		// We can stop if we have the last stream index and this is that index.
		if idx.LastStreamIndex > 0 {
			// Get the index for this entry.
			//
			// We can ignore this error, since "Get" will have already resolved the
			// index successfully.
			if sidx, _ := e.GetStreamIndex(); sidx == types.MessageIndex(idx.LastStreamIndex) {
				return false
			}
		}
		return true
	})
	switch {
	case err != nil:
		return nil, err

	case lastEntry == nil:
		return nil, storage.ErrDoesNotExist

	default:
		return lastEntry, nil
	}
}

// getIndex returns the cached log stream index, fetching it if necessary.
func (s *storageImpl) getIndex() (*logpb.LogIndex, error) {
	s.indexMu.Lock()
	defer s.indexMu.Unlock()

	if s.index == nil {
		index, err := loadIndex(s, s.Client, s.indexPath)
		switch errors.Unwrap(err) {
		case nil:
			break

		case cloudStorage.ErrBucketNotExist, cloudStorage.ErrObjectNotExist:
			// Treat a missing index the same as an empty index.
			log.WithError(err).Warningf(s, "Index is invalid, using empty index.")
			index = &logpb.LogIndex{}

		default:
			return nil, err
		}

		s.index = index
	}
	return s.index, nil
}

func loadIndex(c context.Context, client gs.Client, path gs.Path) (*logpb.LogIndex, error) {
	// If there is no path, then return an empty index.
	if path == "" {
		log.Infof(c, "No index path, using empty index.")
		return &logpb.LogIndex{}, nil
	}

	r, err := client.NewReader(path, 0, -1)
	if err != nil {
		log.WithError(err).Errorf(c, "Failed to create index Reader.")
		return nil, errors.Annotate(err).Reason("failed to create index Reader").Err()
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.WithError(err).Warningf(c, "Error closing index Reader.")
		}
	}()
	indexData, err := ioutil.ReadAll(r)
	if err != nil {
		log.WithError(err).Errorf(c, "Failed to read index.")
		return nil, errors.Annotate(err).Reason("failed to read index").Err()
	}

	index := logpb.LogIndex{}
	if err := proto.Unmarshal(indexData, &index); err != nil {
		log.WithError(err).Errorf(c, "Failed to unmarshal index.")
		return nil, errors.Annotate(err).Reason("failed to unmarshal index").Err()
	}

	return &index, nil
}

type getStrategy struct {
	// startIndex is desired initial log entry index.
	startIndex types.MessageIndex

	// startOffset is the beginning byte offset of the log entry stream. This may
	// be lower than the offset of the starting record if the index is sparse.
	startOffset uint64
	// endOffset is the ending byte offset of the log entry stream. This will be
	// 0 if an end offset is not known.
	endOffset uint64

	// count is the number of log entries that will be fetched. If 0, no upper
	// bound was calculated.
	count uint64
}

func (gs *getStrategy) length() int64 {
	if gs.startOffset < gs.endOffset {
		return int64(gs.endOffset - gs.startOffset)
	}
	return -1
}

// setCount sets the `count` field. If called multiple times, the smallest
// assigned value will be retained.
func (gs *getStrategy) setCount(v uint64) {
	if gs.count == 0 || gs.count > v {
		gs.count = v
	}
}

func buildGetStrategy(req *storage.GetRequest, idx *logpb.LogIndex) *getStrategy {
	st := getStrategy{
		startIndex: req.Index,
	}

	// If the user has requested an index past the end of the stream, return no
	// entries (count == 0). This only works if the last stream index is known.
	if idx.LastStreamIndex > 0 && req.Index > types.MessageIndex(idx.LastStreamIndex) {
		return nil
	}

	// Identify the closest index entry to the requested log.
	//
	// If the requested log starts before the first index entry, we must read from
	// record #0.
	startIndexEntry := indexEntryFor(idx.Entries, req.Index)
	if startIndexEntry >= 0 {
		st.startOffset = idx.Entries[startIndexEntry].Offset
	}

	// Determine an upper bound based on our limits.
	//
	// If we have a count limit, identify the maximum entry that can be loaded,
	// find the index entry closest to it, and use that to determine our upper
	// bound.
	if req.Limit > 0 {
		st.setCount(uint64(req.Limit))

		// Find the index entry for the stream entry AFTER the last one we are going
		// to return.
		entryAfterGetBlock := req.Index + types.MessageIndex(req.Limit)
		endIndexEntry := indexEntryFor(idx.Entries, entryAfterGetBlock)
		switch {
		case endIndexEntry < 0:
			// The last possible request log entry is before the first index entry.
			// Read up to the first index entry.
			endIndexEntry = 0

		case endIndexEntry <= startIndexEntry:
			// The last possible request log entry is closest to the start index
			// entry. Use the index entry immediately after it.
			endIndexEntry = startIndexEntry + 1

		default:
			// We have the index entry <= the stream entry after the last one that we
			// will return.
			//
			// If we're sparse, this could be the index at or before our last entry.
			// If this is the case, use the next index entry, which will be after
			// "entryAfterGetBlock" (EAGB).
			//
			// START ------ LIMIT     (LIMIT+1)
			//   |          [IDX]         |          [IDX]
			// index          |   entryAfterGetBlock   |
			//           endIndexEntry         (endIndexEntry+1)
			if types.MessageIndex(idx.Entries[endIndexEntry].StreamIndex) < entryAfterGetBlock {
				endIndexEntry++
			}
		}

		// If we're pointing to a valid index entry, set our upper bound.
		if endIndexEntry < len(idx.Entries) {
			st.endOffset = idx.Entries[endIndexEntry].Offset
		}
	}

	return &st
}

// indexEntryFor identifies the log index entry closest (<=) to the specified
// index.
//
// If the first index entry is greater than our search index, -1 will be
// returned. This should never happen in practice, though, since our index
// construction always indexes log entry #0.
//
// It does this by performing a binary search over the index entries.
func indexEntryFor(entries []*logpb.LogIndex_Entry, i types.MessageIndex) int {
	ui := uint64(i)
	s := sort.Search(len(entries), func(i int) bool {
		return entries[i].StreamIndex > ui
	})

	// The returned index is the one immediately after the index that we want. If
	// our search returned 0, the first index entry is > our search entry, and we
	// will return nil.
	return s - 1
}
