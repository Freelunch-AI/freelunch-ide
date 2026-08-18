// Package buildinfo exposes the version metadata stamped into the freelunch
// binary at link time.
//
// Release builds set these via -ldflags (see .goreleaser.yaml). Builds produced
// any other way fall back to the values embedded by the Go toolchain from the
// enclosing VCS checkout, so a `go build` from a dirty working tree still
// reports something truthful rather than claiming to be a release.
package buildinfo

import (
	"runtime"
	"runtime/debug"
	"sync"
)

// Injected via -ldflags at release time. Do not read these directly; use Get.
var (
	version = ""
	commit  = ""
	date    = ""
)

// DevVersion is reported when no release version was stamped in.
const DevVersion = "0.0.0-dev"

// Info describes the provenance of a freelunch binary.
type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Date      string `json:"date"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

var (
	once   sync.Once
	cached Info
)

// Get returns the build metadata for the running binary. The result is computed
// once and reused.
func Get() Info {
	once.Do(func() { cached = resolve(debug.ReadBuildInfo) })
	return cached
}

// resolve builds an Info from the linker-injected values, falling back to the
// toolchain's VCS stamps. It takes the reader as an argument so tests can
// exercise the fallback path without a real build.
func resolve(read func() (*debug.BuildInfo, bool)) Info {
	info := Info{
		Version:   version,
		Commit:    commit,
		Date:      date,
		GoVersion: runtime.Version(),
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}

	bi, ok := read()
	if ok {
		settings := make(map[string]string, len(bi.Settings))
		for _, s := range bi.Settings {
			settings[s.Key] = s.Value
		}
		if info.Commit == "" {
			info.Commit = settings["vcs.revision"]
			if settings["vcs.modified"] == "true" && info.Commit != "" {
				info.Commit += "-dirty"
			}
		}
		if info.Date == "" {
			info.Date = settings["vcs.time"]
		}
		if info.Version == "" && bi.Main.Version != "" && bi.Main.Version != "(devel)" {
			info.Version = bi.Main.Version
		}
	}

	if info.Version == "" {
		info.Version = DevVersion
	}
	if info.Commit == "" {
		info.Commit = "unknown"
	}
	if info.Date == "" {
		info.Date = "unknown"
	}

	return info
}
