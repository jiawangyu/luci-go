// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package pubsub

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/common/logging/gologger"

	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/clock/testclock"
	"github.com/luci/luci-go/common/data/recordio"
	gcps "github.com/luci/luci-go/common/gcloud/pubsub"
	"github.com/luci/luci-go/common/proto/google"
	"github.com/luci/luci-go/grpc/grpcutil"
	"github.com/luci/luci-go/logdog/api/logpb"

	"cloud.google.com/go/pubsub"
	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

type testTopic struct {
	sync.Mutex

	err func() error

	msgC          chan *pubsub.Message
	nextMessageID int
}

func (t *testTopic) String() string { return "test" }

func (t *testTopic) Publish(c context.Context, msgs ...*pubsub.Message) ([]string, error) {
	if t.err != nil {
		if err := t.err(); err != nil {
			return nil, err
		}
	}

	ids := make([]string, len(msgs))
	for i, m := range msgs {
		if t.msgC != nil {
			select {
			case t.msgC <- m:
			case <-c.Done():
				return nil, c.Err()
			}
		}
		ids[i] = t.getNextMessageID()
	}
	return ids, nil
}

func (t *testTopic) getNextMessageID() string {
	t.Lock()
	defer t.Unlock()

	id := t.nextMessageID
	t.nextMessageID++
	return fmt.Sprintf("%d", id)
}

func deconstructMessage(msg *pubsub.Message) (*logpb.ButlerMetadata, *logpb.ButlerLogBundle, error) {
	fr := recordio.NewReader(bytes.NewBuffer(msg.Data), gcps.MaxPublishSize)

	// Validate header frame.
	headerBytes, err := fr.ReadFrameAll()
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to read header frame: %s", err)
	}

	header := logpb.ButlerMetadata{}
	if err := proto.Unmarshal(headerBytes, &header); err != nil {
		return nil, nil, fmt.Errorf("test: failed to unmarshal header: %s", err)
	}

	if header.Type != logpb.ButlerMetadata_ButlerLogBundle {
		return nil, nil, fmt.Errorf("test: unknown frame data type: %v", header.Type)
	}

	// Validate data frame.
	data, err := fr.ReadFrameAll()
	if err != nil {
		return nil, nil, fmt.Errorf("test: failed to read data frame: %s", err)
	}

	switch header.Compression {
	case logpb.ButlerMetadata_ZLIB:
		r, err := zlib.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, nil, fmt.Errorf("test: failed to create zlib reader: %s", err)
		}
		defer r.Close()

		data, err = ioutil.ReadAll(r)
		if err != nil {
			return nil, nil, fmt.Errorf("test: failed to read compressed data: %s", err)
		}
	}

	dataBundle := logpb.ButlerLogBundle{}
	if err := proto.Unmarshal(data, &dataBundle); err != nil {
		return nil, nil, fmt.Errorf("test: failed to unmarshal bundle: %s", err)
	}

	return &header, &dataBundle, nil
}

func TestOutput(t *testing.T) {
	Convey(`An Output using a test Pub/Sub instance`, t, func() {
		ctx, tc := testclock.UseTime(context.Background(), time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC))
		ctx = gologger.StdConfig.Use(ctx)
		ctx = logging.SetLevel(ctx, logging.Debug)
		tt := &testTopic{
			msgC: make(chan *pubsub.Message),
		}
		conf := Config{
			Topic: tt,
		}
		o := New(ctx, conf).(*pubSubOutput)
		So(o, ShouldNotBeNil)
		defer o.Close()

		bundle := &logpb.ButlerLogBundle{
			Timestamp: google.NewTimestamp(clock.Now(ctx)),
			Entries: []*logpb.ButlerLogBundle_Entry{
				{},
			},
		}

		Convey(`Can send/receive a bundle.`, func() {
			errC := make(chan error)
			go func() {
				errC <- o.SendBundle(bundle)
			}()
			msg := <-tt.msgC
			So(<-errC, ShouldBeNil)

			h, b, err := deconstructMessage(msg)
			So(err, ShouldBeNil)
			So(h.Compression, ShouldEqual, logpb.ButlerMetadata_NONE)
			So(b, ShouldResemble, bundle)

			Convey(`And records stats.`, func() {
				st := o.Stats()
				So(st.Errors(), ShouldEqual, 0)
				So(st.SentBytes(), ShouldBeGreaterThan, 0)
				So(st.SentMessages(), ShouldEqual, 1)
				So(st.DiscardedMessages(), ShouldEqual, 0)
			})
		})

		Convey(`Will return an error if Publish failed non-transiently.`, func() {
			tt.err = func() error { return grpcutil.InvalidArgument }
			So(o.SendBundle(bundle), ShouldEqual, grpcutil.InvalidArgument)
		})

		Convey(`Will retry indefinitely if Publish fails transiently (Context deadline).`, func() {
			const retries = 30

			// Advance our clock each time there is a delay up until count.
			count := 0
			tc.SetTimerCallback(func(d time.Duration, t clock.Timer) {
				switch {
				case !testclock.HasTags(t, clock.ContextDeadlineTag):
					// Other timer (probably retry sleep), advance time.
					tc.Add(d)

				case count < retries:
					// Still retrying, advance time.
					count++
					tc.Add(d)

				default:
					// Done retrying, don't expire Contexts anymore. Consume the message
					// when it is sent.
					tt.err = func() error { return grpcutil.InvalidArgument }
				}
			})

			// Time our our RPC. Because of our timer callback, this will always be
			// hit.
			o.RPCTimeout = 30 * time.Second
			So(o.SendBundle(bundle), ShouldEqual, grpcutil.InvalidArgument)
			So(count, ShouldEqual, retries)
		})

		Convey(`Will retry indefinitely if Publish fails transiently (gRPC).`, func() {
			const retries = 30

			// Advance our clock each time there is a delay up until count.
			tc.SetTimerCallback(func(d time.Duration, t clock.Timer) {
				tc.Add(d)
			})

			count := 0
			tt.msgC = nil
			tt.err = func() error {
				count++
				if count < retries {
					return grpcutil.Internal
				}
				return grpcutil.NotFound // Non-transient.
			}

			// Time our our RPC. Because of our timer callback, this will always be
			// hit.
			So(o.SendBundle(bundle), ShouldEqual, grpcutil.NotFound)
			So(count, ShouldEqual, retries)
		})
	})
}
