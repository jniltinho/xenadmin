// Package buildinfo holds the version stamped at build time (-ldflags "-X xcpng-center/internal/buildinfo.Version=...").
package buildinfo

var (
	Version = "dev"
	Commit  = ""
)
