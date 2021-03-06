// Copyright 2017 The LUCI Authors.
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

package rpc

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/grpc/grpcutil"
	"go.chromium.org/luci/server/auth"

	"go.chromium.org/luci/machine-db/api/crimson/v1"
)

// authPrelude ensures the user is authorized to use the Crimson API.
func authPrelude(c context.Context, methodName string, req proto.Message) (context.Context, error) {
	groups := []string{"machine-db-administrators"}
	if strings.HasPrefix(methodName, "List") {
		groups = append(groups, "machine-db-readers")
	}
	switch is, err := auth.IsMember(c, groups...); {
	case err != nil:
		return c, err
	case !is:
		return c, status.Errorf(codes.PermissionDenied, "unauthorized user")
	}
	return c, nil
}

// gRPCifyAndLogErr ensures any error being returned is a gRPC error, logging Internal and Unknown errors.
func gRPCifyAndLogErr(c context.Context, methodName string, rsp proto.Message, err error) error {
	return grpcutil.GRPCifyAndLogErr(c, err)
}

// matches returns whether the given string matches the given set.
// An empty set matches all strings.
func matches(s string, set stringset.Set) bool {
	return set.Has(s) || set.Len() == 0
}

// NewServer returns a new Crimson RPC server.
func NewServer() crimson.CrimsonServer {
	return &crimson.DecoratedCrimson{
		Prelude:  authPrelude,
		Service:  &Service{},
		Postlude: gRPCifyAndLogErr,
	}
}

// Service handles Crimson RPCs.
// Not a proper gRPC service, errors may not be gRPC errors. Use NewServer to get a proper gRPC service.
type Service struct {
}
