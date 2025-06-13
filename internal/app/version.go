package app

import (
	"fmt"
	"runtime"
)

var (
	Version   = "v0.1.0"
	Commit    = "unknown"
	BuildTime = "unknown"
)

func VersionInfo() string {
	return fmt.Sprintf("hass-cli %s (%s) built at %s with %s",
		Version, Commit, BuildTime, runtime.Version())
}