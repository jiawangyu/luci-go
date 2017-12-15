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

package config

import (
	"golang.org/x/net/context"

	"go.chromium.org/luci/common/config/validation"
	"go.chromium.org/luci/common/data/stringset"
	"go.chromium.org/luci/common/errors"
	"go.chromium.org/luci/common/logging"
	"go.chromium.org/luci/luci_config/common/cfgtypes"
	"go.chromium.org/luci/luci_config/server/cfgclient"
	"go.chromium.org/luci/luci_config/server/cfgclient/textproto"

	"go.chromium.org/luci/machine-db/api/config/v1"
)

// platformsConfig is the name of the config file enumerating platforms.
const platformsConfig = "platforms.cfg"

// importPlatformConfigs fetches, validates, and applies platform configs.
func importPlatformConfigs(c context.Context, configSet cfgtypes.ConfigSet) error {
	platform := &config.PlatformsConfig{}
	metadata := &cfgclient.Meta{}
	if err := cfgclient.Get(c, cfgclient.AsService, configSet, platformsConfig, textproto.Message(platform), metadata); err != nil {
		return errors.Annotate(err, "failed to load %s", platformsConfig).Err()
	}
	logging.Infof(c, "Found %s revision %q", platformsConfig, metadata.Revision)

	validationContext := &validation.Context{Context: c}
	validationContext.SetFile(platformsConfig)
	validatePlatformsConfig(validationContext, platform)
	if err := validationContext.Finalize(); err != nil {
		return errors.Annotate(err, "invalid config").Err()
	}
	return nil
}

// validatePlatformsConfig validates platforms.cfg.
func validatePlatformsConfig(c *validation.Context, cfg *config.PlatformsConfig) {
	// Platform names must be unique.
	// Keep records of ones we've already seen.
	names := stringset.New(len(cfg.Platform))

	for _, p := range cfg.Platform {
		switch {
		case p.Name == "":
			c.Errorf("platform names are required and must be non-empty")
		case !names.Add(p.Name):
			c.Errorf("duplicate platform %q", p.Name)
		}
	}
}
