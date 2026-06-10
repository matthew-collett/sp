package build

import (
	"runtime/debug"

	"github.com/matthew-collett/sp/internal/ui"
)

type VersionInfo struct {
	Version string
	Date    string
}

var (
	v       = "dev"
	d       = ""
	Version = &VersionInfo{Version: v, Date: d}
)

func init() {
	// Fallback to runtime/debug build info.
	if v == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			v = info.Main.Version
		}
	}
	Version.Version = v
	Version.Date = d
}

func (vi *VersionInfo) Style() *ui.Style {
	if vi.Date != "" {
		return ui.Success("sp %s (%s)", vi.Version, vi.Date)
	}
	return ui.Success("sp %s", vi.Version)
}

func (vi *VersionInfo) Template() string {
	return vi.Style().String() + "\n"
}

func (vi *VersionInfo) String() string {
	return vi.Version
}
