package manifests

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/sneha-afk/trovl/internal/links"
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

var allSupportedPlatforms []string = []string{
	"windows",
	"linux",
	"darwin",
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

func (m *Manifest) Apply(verbose bool, constructOpts *links.ConstructOptions) error {
	for _, manifestLink := range m.Links {
		// 1. Separate platform using the top-level (default) symlink vs. overrides
		platformUsingSpec := mapset.NewSet[string]()

		if slices.Contains(manifestLink.Platforms, "all") {
			platformUsingSpec.Append(allSupportedPlatforms...)
		} else {
			platformUsingSpec.Append(manifestLink.Platforms...)
		}

		keys := make([]string, 0, len(manifestLink.PlatformOverrides))
		for k := range manifestLink.PlatformOverrides {
			keys = append(keys, k)
		}
		platformUsingSpec.RemoveAll(keys...)

		// 2. Detect current OS and carry out links
		var linkToUse string
		if platformUsingSpec.Contains(runtime.GOOS) {

			linkToUse = manifestLink.Link
		} else {
			linkToUse = manifestLink.PlatformOverrides[runtime.GOOS].Link
		}

		// TODO: allow granular ovewrites?
		linkSpec, err := links.Construct(manifestLink.Target, linkToUse, manifestLink.Relative, constructOpts)
		if errors.Is(err, links.ErrDeclinedOverwrite) {
			continue
		}
		if err != nil {
			return fmt.Errorf("could not construct link: %v", err)
		}
		if err := links.Add(linkSpec); err != nil {
			return fmt.Errorf("could not add link: %v", err)
		}

		// TODO: verbosity here
	}
	return nil
}
