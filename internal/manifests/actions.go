/*
Package manifests defines a link manifest that is used for bulk operations.
Pretty much the exact reason why this project exists.
*/
package manifests

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sneha-afk/trovl/internal/links"
	"github.com/sneha-afk/trovl/internal/models"
	"github.com/sneha-afk/trovl/internal/state"
)

type PlatformOverride struct {
	Link string `json:"link"`
}

type ManifestLink struct {
	Target            string                      `json:"target"`
	Link              string                      `json:"link"`
	Platforms         []string                    `json:"platforms"`
	Relative          bool                        `json:"relative"`
	PlatformOverrides map[string]PlatformOverride `json:"platform_overrides,omitempty"`
}

type Manifest struct {
	Links []ManifestLink `json:"links"`
}

var allSupportedPlatforms mapset.Set[string] = mapset.NewSet("windows", "linux", "darwin")

func IsSupportedPlatform(platform string) bool {
	platform = strings.ToLower(platform)
	if platform == "all" {
		return true
	}
	return allSupportedPlatforms.Contains(platform)
}

func New(path string) (*Manifest, error) {
	manifestFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read manifest file: %v", err)
	}

	m := &Manifest{}
	if err := json.Unmarshal(manifestFile, &m); err != nil {
		return nil, fmt.Errorf("could not unmarshal manifest: %v", err)
	}
	m.FillDefaults()

	return m, nil
}

func (m *Manifest) FillDefaults() {
	for i := range m.Links {
		if len(m.Links[i].Platforms) == 0 {
			m.Links[i].Platforms = append(m.Links[i].Platforms, "all")
		}
	}
}

func (m *Manifest) UnmarshalJSON(data []byte) error {
	// alias ensures no infinte loop since the alias has nothing defined on it
	type manifestAlias Manifest

	var temp manifestAlias
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	for i := range temp.Links {
		link := &temp.Links[i]

		if len(link.Platforms) == 0 {
			link.Platforms = []string{"all"}
		}

		if link.PlatformOverrides == nil {
			link.PlatformOverrides = map[string]PlatformOverride{}
		}

		if link.Target == "" {
			return fmt.Errorf("links[%d]: missing target", i)
		}
		if link.Link == "" {
			return fmt.Errorf("links[%d]: missing link", i)
		}

		if slices.Contains(link.Platforms, "all") && len(link.Platforms) > 1 {
			return fmt.Errorf("links[%d]: 'all' cannot be combined with other platforms", i)
		}

		seen := map[string]struct{}{}
		for _, plat := range link.Platforms {
			if !IsSupportedPlatform(plat) {
				return fmt.Errorf("links[%d]: unsupported platform %q", i, plat)
			}
			if _, ok := seen[plat]; ok {
				return fmt.Errorf("links[%d]: duplicate platform %q", i, plat)
			}
			seen[plat] = struct{}{}
		}

		for plat, over := range link.PlatformOverrides {
			if !IsSupportedPlatform(plat) {
				return fmt.Errorf("links[%d] (override): unsupported platform %q", i, plat)
			}
			if over.Link == "" {
				return fmt.Errorf("links[%d]: override %q missing link", i, plat)
			}
		}
	}

	*m = Manifest(temp)
	return nil
}

func (m *Manifest) Apply(state *state.TrovlState) error {
	var constructed models.Link
	var err error
	var linkToUse string
	var numLinks = len(m.Links)

	for i := range m.Links {
		link := &m.Links[i]

		if override, ok := link.PlatformOverrides[runtime.GOOS]; ok {
			// 1. If an override exists for this platform, it always wins
			linkToUse = override.Link
		} else {
			// 2. Determine whether this link applies to the current platform
			if slices.Contains(link.Platforms, "all") || slices.Contains(link.Platforms, runtime.GOOS) {
				linkToUse = link.Link
			} else {
				return fmt.Errorf("links[%d]: link does not apply to platform %q", i, runtime.GOOS)
			}
		}

		constructed, err = links.Construct(state, link.Target, linkToUse)
		if err != nil && !errors.Is(err, links.ErrDeclinedOverwrite) {
			return fmt.Errorf("links[%d]: %w", i, err)
		}
		if err := links.Add(&constructed); err != nil {
			return fmt.Errorf("links[i] %w", err)
		}

		if state.Verbose() {
			state.Logger.Info(fmt.Sprintf("Successfully added symlink [%v/%v]", i, numLinks), "link", linkToUse, "target", link.Target)
		}
	}

	return nil
}
