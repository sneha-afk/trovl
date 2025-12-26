package models

type PlatformOverride struct {
	Platform string `json:"platform"`
	Link     string `json:"link"`
}

type ManifestLink struct {
	Target            string             `json:"target"`
	Link              string             `json:"link"`
	Platforms         []string           `json:"platforms"`
	Relative          bool               `json:"relative"`
	PlatformOverrides []PlatformOverride `json:"platform_overrides,omitempty"`
}

type Manifest struct {
	Links []ManifestLink `json:"links"`
}

func (m *Manifest) FillDefaults() {
	for i := range m.Links {
		if len(m.Links[i].Platforms) == 0 {
			m.Links[i].Platforms = append(m.Links[i].Platforms, "all")
		}
	}
}
