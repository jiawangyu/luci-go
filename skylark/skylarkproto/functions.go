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

package skylarkproto

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
)

// ProtoLib() returns a dict with single struct named "proto" that holds helper
// functions to manipulate protobuf messages (in particular serialize them).
func ProtoLib() skylark.StringDict {
	return skylark.StringDict{
		"proto": skylarkstruct.FromStringDict(skylark.String("proto"), skylark.StringDict{
			"to_pbtext": skylark.NewBuiltin("to_pbtext", toPbText),
		}),
	}
}

// toPbText takes single protobuf message and serializes it using text protobuf
// serialization.
func toPbText(_ *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var msg *Message
	if err := skylark.UnpackArgs("to_pbtext", args, kwargs, "msg", &msg); err != nil {
		return nil, err
	}
	pb, err := msg.ToProto()
	if err != nil {
		return nil, err
	}
	return skylark.String(proto.MarshalTextString(pb)), nil
}
