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

package git

import (
	"net/http"
	"strings"

	"go.chromium.org/luci/common/api/gitiles"
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/errors"
	gitilespb "go.chromium.org/luci/common/proto/gitiles"
	"go.chromium.org/luci/server/auth"
	"golang.org/x/net/context"
)

// ClientFactory creates a Gitiles client.
type ClientFactory func(ctx context.Context, host string) (gitilespb.GitilesClient, error)

var factoryKey = "gitiles client factory key"

// UseFactory installs f into c.
func UseFactory(c context.Context, f ClientFactory) context.Context {
	return context.WithValue(c, &factoryKey, f)
}

// TODO(tandrii): remove the following per https://crbug.com/796317.
// Until Milo properly supports ACLs for blamelists, we have a hack; if the git
// repo being log'd is in this list, use `auth.AsSelf`. Otherwise use
// `auth.Anonymous`.
//
// The reason to do this is that we currently do blamelist calculation in the
// backend, so we can't accurately determine if the requesting user has access
// to these repos or not. For now, we use this whitelist to indicate domains
// that we know have full public read-access so that we can use milo's
// credentials (instead of anonymous) in order to avoid hitting gitiles'
// anonymous quota limits.
var whitelistPublicHosts = stringset.NewFromSlice(
	"chromium",
	"fuchsia",
)

func isPublicHost(host string) bool {
	const gs = ".googlesource.com"
	if !strings.HasSuffix(host, gs) {
		return false
	}

	host = strings.TrimSuffix(host, gs)
	host = strings.TrimSuffix(host, "-review")
	return whitelistPublicHosts.Has(host)
}

// Transport returns an HTTP transport to be used for the given host.
//
// Currently, it is authenticated as self only for a whitelistPublicHosts.
// For all other repos, the transport will not use authentication to avoid
// information leaks.
// TODO(tandrii): fix this per https://crbug.com/796317.
func Transport(c context.Context, host string) (transport http.RoundTripper, authenticated bool, err error) {
	if isPublicHost(host) {
		transport, err = auth.GetRPCTransport(c, auth.AsSelf, auth.WithScopes(gitiles.OAuthScope))
		authenticated = true
	} else {
		transport, err = auth.GetRPCTransport(c, auth.NoAuth)
		authenticated = false
	}
	return
}

// GitilesProdClient returns a production Gitiles client.
// Implements ClientFactory.
func GitilesProdClient(c context.Context, host string) (gitilespb.GitilesClient, error) {
	t, auth, err := Transport(c, host)
	if err != nil {
		return nil, errors.Annotate(err, "getting RPC Transport").Err()
	}
	return gitiles.NewRESTClient(&http.Client{Transport: t}, host, auth)
}

// Client creates a new Gitiles client using the ClientFactory installed in c.
// See also UseFactory.
func Client(c context.Context, host string) (gitilespb.GitilesClient, error) {
	f, ok := c.Value(&factoryKey).(ClientFactory)
	if !ok {
		return nil, errors.New("gitiles client factory is not installed in context")
	}
	return f(c, host)
}
