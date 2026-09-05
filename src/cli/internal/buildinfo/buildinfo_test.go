package buildinfo

import (
	"runtime/debug"
	"testing"
)

func TestResolveFallsBackToDevVersion(t *testing.T) {
	got := resolve(func() (*debug.BuildInfo, bool) { return nil, false })

	if got.Version != DevVersion {
		t.Errorf("Version = %q, want %q", got.Version, DevVersion)
	}
	if got.Commit != "unknown" {
		t.Errorf("Commit = %q, want %q", got.Commit, "unknown")
	}
	if got.Date != "unknown" {
		t.Errorf("Date = %q, want %q", got.Date, "unknown")
	}
	if got.GoVersion == "" || got.Platform == "" {
		t.Errorf("GoVersion/Platform must always be populated, got %+v", got)
	}
}

func TestResolveUsesVCSStampsWhenNotStamped(t *testing.T) {
	got := resolve(func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Main: debug.Module{Version: "(devel)"},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "abc123"},
				{Key: "vcs.modified", Value: "true"},
				{Key: "vcs.time", Value: "2026-07-27T00:00:00Z"},
			},
		}, true
	})

	if got.Commit != "abc123-dirty" {
		t.Errorf("Commit = %q, want %q", got.Commit, "abc123-dirty")
	}
	if got.Date != "2026-07-27T00:00:00Z" {
		t.Errorf("Date = %q, want %q", got.Date, "2026-07-27T00:00:00Z")
	}
	if got.Version != DevVersion {
		t.Errorf("Version = %q, want %q for a (devel) main module", got.Version, DevVersion)
	}
}

func TestResolvePrefersLinkerValues(t *testing.T) {
	t.Cleanup(func() { version, commit, date = "", "", "" })
	version, commit, date = "1.2.3", "deadbeef", "2026-01-01T00:00:00Z"

	got := resolve(func() (*debug.BuildInfo, bool) {
		return &debug.BuildInfo{
			Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "ignored"}},
		}, true
	})

	if got.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", got.Version, "1.2.3")
	}
	if got.Commit != "deadbeef" {
		t.Errorf("Commit = %q, want linker value %q", got.Commit, "deadbeef")
	}
}
