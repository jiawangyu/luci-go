// Copyright 2018 The LUCI Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metadata

import (
	"golang.org/x/net/context"

	api "go.chromium.org/luci/cipd/api/cipd/v1"
)

// Storage knows how to store, fetch and update prefix metadata, as well as
// how to calculate its fingerprint.
//
// Doesn't try to understand what metadata means, just fingerprints and stores
// it.
//
// This functionality is organized into an interface to simplify mocking. Use
// GetStorage to grab a real implementation.
type Storage interface {
	// GetMetadata fetches metadata associated with the given prefix and all
	// parent prefixes.
	//
	// Does not check permissions.
	//
	// The return value is sorted by the prefix length. Prefixes without metadata
	// are skipped. For example, when requesting metadata for prefix "a/b/c/d" the
	// return value may contain entries for "a", "a/b", "a/b/c/d" (in that order,
	// with "a/b/c" skipped in this example as not having any metadata attached).
	//
	// Note that the prefix of the last entry doesn't necessary match 'prefix'.
	// This happens if metadata for that prefix doesn't exist. Similarly, the
	// returns value may be completely empty slice in case there's no metadata
	// for the requested prefix and all its parent prefixes.
	//
	// Returns a fatal error if the prefix is malformed, all other errors are
	// transient.
	GetMetadata(c context.Context, prefix string) ([]*api.PrefixMetadata, error)

	// UpdateMetadata transactionally updates or creates metadata of some
	// prefix.
	//
	// Does not check permissions. Does not check the format of the updated
	// metadata.
	//
	// If fetches the metadata object and calls the callback to modify it. The
	// callback may be called multiple times when retrying the transaction. If the
	// callback doesn't return an error, the new metadata's fingerprint is updated
	// and the metadata is saved to the storage and returned to the caller.
	//
	// If the metadata object doesn't exist yet, the callback will be called with
	// an empty object that has only 'Prefix' field populated. The callback then
	// can populate the rest of the fields.
	//
	// If the callback returns an error, it will be returned as is. If the
	// transaction itself fails, returns a transient error.
	UpdateMetadata(c context.Context, prefix string, cb func(m *api.PrefixMetadata) error) (*api.PrefixMetadata, error)
}

// GetStorage returns production implementation of the metadata storage.
func GetStorage() Storage {
	return &legacyStorageImpl{}
}