package models

type LinkType int

const (
	LinkFile LinkType = iota
	LinkDirectory
	// TODO: differentiate junctions on Windows?
)

type Link struct {
	Target    string   `json:"target"`     // Real file/directory
	LinkMount string   `json:"link_mount"` // Where the symlink is
	Type      LinkType `json:"link_type"`
}
