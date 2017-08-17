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

package ensure

import (
	"fmt"
	"net/url"

	"go.chromium.org/luci/cipd/client/cipd/common"
	"go.chromium.org/luci/common/errors"
)

// an itemParser should parse the value from `val`, and update s or
// f accordingly, returning an error if needed.
type itemParser func(s *itemParserState, f *File, val string) error

// itemParserState is the state object shared between the item parsers and the
// main ParseFile implementation.
type itemParserState struct {
	curSubdir string
}

func subdirParser(s *itemParserState, _ *File, val string) (err error) {
	// We expand with the default expander here just to see if this is a plausible
	// template. When the user uses File.ResolveWith, this will actually use the
	// user-supplied expander.
	tempExpanded := ""
	if tempExpanded, err = common.DefaultTemplateExpander().Validate(val); err == nil {
		if err = common.ValidateSubdir(tempExpanded); err == nil {
			s.curSubdir = val
		}
	} else {
		err = errors.Annotate(err, "bad subdir %q", val).Err()
	}
	return
}

func serviceURLParser(_ *itemParserState, f *File, val string) error {
	if f.ServiceURL != "" {
		return fmt.Errorf("$ServiceURL may only be set once per file")
	}
	if _, err := url.Parse(val); err != nil {
		return fmt.Errorf("expecting '$ServiceURL <url>' but url is invalid: %s", err)
	}
	f.ServiceURL = val
	return nil
}

// itemParsers is the main way that the ensure file format is extended. If you
// need to add a new setting or directive, please add an appropriate function
// above and then add it to this map.
var itemParsers = map[string]itemParser{
	"@subdir":     subdirParser,
	"$serviceurl": serviceURLParser,
}
