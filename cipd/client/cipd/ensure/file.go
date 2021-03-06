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
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/url"
	"sort"
	"strings"
	"unicode"

	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/iotools"

	"go.chromium.org/luci/cipd/client/cipd/template"
	"go.chromium.org/luci/cipd/common"
)

// File is an in-process representation of the 'ensure file' format.
type File struct {
	ServiceURL string

	PackagesBySubdir map[string]PackageSlice
	VerifyPlatforms  []template.Platform
}

// ParseFile parses an ensure file from the given reader. See the package docs
// for the format of this file.
//
// This file will contain unresolved template strings for package names as well
// as unpinned package versions. Use File.Resolve() to obtain resolved+pinned
// versions of these.
func ParseFile(r io.Reader) (*File, error) {
	ret := &File{PackagesBySubdir: map[string]PackageSlice{}}

	state := itemParserState{}

	// indicates that the parser is able to read $setting lines. This is flipped
	// to false on the first non-$setting statement in the file.
	settingsAllowed := true

	lineNo := 0
	makeError := func(fmtStr string, args ...interface{}) error {
		args = append([]interface{}{lineNo}, args...)
		return fmt.Errorf("failed to parse desired state (line %d): "+fmtStr, args...)
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lineNo++

		// Remove all space
		line := strings.TrimSpace(scanner.Text())

		if len(line) == 0 || line[0] == '#' {
			// skip blank lines and comments
			continue
		}

		tok1 := line
		tok2 := ""
		if idx := strings.IndexFunc(line, unicode.IsSpace); idx == -1 {
			// only one token. This implies a second token of ""
		} else {
			tok1, tok2 = line[:idx], strings.TrimSpace(line[idx:])
		}

		switch c := tok1[0]; c {
		case '@', '$':
			if c == '$' {
				if !settingsAllowed {
					return nil, makeError("$setting found after non-$setting statements")
				}
			} else {
				settingsAllowed = false
			}

			if p := itemParsers[strings.ToLower(tok1)]; p != nil {
				if err := p(&state, ret, tok2); err != nil {
					return nil, makeError("%s", err)
				}
			} else {
				tag := map[byte]string{'@': "@directive", '$': "$setting"}[c]
				return nil, makeError("unknown %s: %q", tag, tok1)
			}

		default:
			settingsAllowed = false
			pkg := PackageDef{tok1, tok2, lineNo}

			_, err := pkg.Resolve(func(pkg, vers string) (common.Pin, error) {
				return common.Pin{
					PackageName: pkg,
					InstanceID:  vers,
				}, common.ValidateInstanceVersion(vers)
			}, template.DefaultExpander())
			if err != nil && err != template.ErrSkipTemplate {
				return nil, err
			}

			ret.PackagesBySubdir[state.curSubdir] = append(ret.PackagesBySubdir[state.curSubdir], pkg)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

// ResolvedFile only contains valid, fully-resolved information and is the
// result of calling File.Resolve.
type ResolvedFile struct {
	ServiceURL string

	PackagesBySubdir common.PinSliceBySubdir
}

// Serialize writes the ResolvedFile to an io.Writer in canonical order.
func (f *ResolvedFile) Serialize(w io.Writer) (int, error) {
	// piggyback on top of File.Serialize.
	packagesBySubdir := make(map[string]PackageSlice, len(f.PackagesBySubdir))
	for k, v := range f.PackagesBySubdir {
		slc := make(PackageSlice, len(v))
		for i, pkg := range v {
			slc[i] = PackageDef{
				PackageTemplate:   pkg.PackageName,
				UnresolvedVersion: pkg.InstanceID,
			}
		}
		packagesBySubdir[k] = slc
	}
	return (&File{f.ServiceURL, packagesBySubdir, nil}).Serialize(w)
}

// Resolve takes the current unresolved File and expands all package templates
// using common.DefaultPackageNameExpander(), and also resolves all versions
// with the provided VersionResolver.
func (f *File) Resolve(rslv VersionResolver) (*ResolvedFile, error) {
	return f.ResolveWith(rslv, template.DefaultExpander())
}

// ResolveWith takes the current unresolved File and expands all package
// templates using the provided values of arch and os, and also resolves
// all versions with the provided VersionResolver.
func (f *File) ResolveWith(rslv VersionResolver, expander template.Expander) (*ResolvedFile, error) {
	ret := &ResolvedFile{}

	if f.ServiceURL != "" {
		// double check the url
		if _, err := url.Parse(f.ServiceURL); err != nil {
			return nil, errors.Annotate(err, "bad ServiceURL").Err()
		}
	}

	ret.ServiceURL = f.ServiceURL
	if len(f.PackagesBySubdir) == 0 {
		return ret, nil
	}

	// subdir -> pkg -> orig_lineno
	resolvedPkgDupList := map[string]map[string]int{}

	ret.PackagesBySubdir = common.PinSliceBySubdir{}
	for subdir, pkgs := range f.PackagesBySubdir {
		realSubdir, err := expander.Expand(subdir)
		switch err {
		case template.ErrSkipTemplate:
			continue
		case nil:
		default:
			return nil, errors.Annotate(err, "normalizing %q", subdir).Err()
		}

		// double-check the subdir
		if err := common.ValidateSubdir(realSubdir); err != nil {
			return nil, errors.Annotate(err, "normalizing %q", subdir).Err()
		}
		for _, pkg := range pkgs {
			pin, err := pkg.Resolve(rslv, expander)
			if err == template.ErrSkipTemplate {
				continue
			}
			if err != nil {
				return nil, errors.Annotate(err, "resolving package").Err()
			}

			if origLineNo, ok := resolvedPkgDupList[realSubdir][pin.PackageName]; ok {
				return nil, errors.
					Reason("duplicate package in subdir %q: %q: defined on line %d and %d",
						realSubdir, pin.PackageName, origLineNo, pkg.LineNo).Err()
			}
			if resolvedPkgDupList[realSubdir] == nil {
				resolvedPkgDupList[realSubdir] = map[string]int{}
			}
			resolvedPkgDupList[realSubdir][pin.PackageName] = pkg.LineNo

			ret.PackagesBySubdir[realSubdir] = append(ret.PackagesBySubdir[realSubdir], pin)
		}
	}

	return ret, nil
}

// Serialize writes the File to an io.Writer in canonical order.
func (f *File) Serialize(w io.Writer) (int, error) {
	return iotools.WriteTracker(w, func(w io.Writer) error {
		needsNLs := 0
		maybeAddNL := func() {
			if needsNLs > 0 {
				w.Write(bytes.Repeat([]byte("\n"), needsNLs))
				needsNLs = 0
			}
		}

		if f.ServiceURL != "" {
			maybeAddNL()
			fmt.Fprintf(w, "$ServiceURL %s", f.ServiceURL)
			needsNLs = 2
		}

		if len(f.VerifyPlatforms) > 0 {
			for _, plat := range f.VerifyPlatforms {
				maybeAddNL()

				fmt.Fprintf(w, "$VerifiedPlatform %s", plat.String())
				needsNLs = 1
			}

			needsNLs = 2
		}

		keys := make(sort.StringSlice, 0, len(f.PackagesBySubdir))
		for k := range f.PackagesBySubdir {
			keys = append(keys, k)
		}
		keys.Sort()

		for _, k := range keys {
			maybeAddNL()
			if k != "" {
				fmt.Fprintf(w, "@Subdir %s", k)
				needsNLs = 1
			}

			pkgs := f.PackagesBySubdir[k]
			pkgsSort := make(PackageSlice, len(pkgs))
			maxLength := 0
			for i, pkg := range pkgs {
				pkgsSort[i] = pkg
				if l := len(pkg.PackageTemplate); l > maxLength {
					maxLength = l
				}
			}
			sort.Sort(pkgsSort)

			for _, p := range pkgsSort {
				maybeAddNL()
				fmt.Fprintf(w, "%-*s %s", maxLength+1, p.PackageTemplate, p.UnresolvedVersion)
				needsNLs = 1
			}
			needsNLs++
		}

		// We only ever want to end the file with 1 newline.
		if needsNLs > 0 {
			needsNLs = 1
		}
		maybeAddNL()
		return nil
	})
}
