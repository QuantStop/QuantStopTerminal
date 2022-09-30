package engine

import (
	"fmt"
	"runtime"
	"time"
)

var (
	Copyright       = fmt.Sprintf("Copyright (c) 2021-%d QuantStop.com", time.Now().Year())
	PrereleaseBlurb = "This version is pre-release and is not intended to be used as a production ready system - use at your own risk."
	GitHub          = "GitHub: https://github.com/QuantStop/QuantStopTerminal"
	Issues          = "Issues: https://github.com/QuantStop/QuantStopTerminal/issues"
)

type Version struct {
	Version         string
	BuildTime       string
	Copyright       string
	PrereleaseBlurb string
	GitHub          string
	Issues          string
	IsDaemon        bool
	IsPreRelease    bool
	IsDevMode       bool
}

func CreateDefaultVersion() *Version {
	return &Version{
		Version:         "0.0.1",
		Copyright:       Copyright,
		PrereleaseBlurb: PrereleaseBlurb,
		GitHub:          GitHub,
		Issues:          Issues,
		IsDaemon:        false,
		IsPreRelease:    false,
		IsDevMode:       true,
	}
}

// GetVersionString returns the version string
func (version *Version) GetVersionString(short bool) string {
	versionStr := fmt.Sprintf(
		"QuantstopTerminal v%s %s %s",
		version.Version,
		runtime.GOARCH,
		runtime.Version())

	if version.IsPreRelease {
		versionStr += " pre-release.\n"
		if !short {
			versionStr += version.PrereleaseBlurb + "\n"
		}
	} else {
		versionStr += " release.\n"
	}

	if version.IsDevMode {
		versionStr += "Development mode: On\n"
	} else {
		versionStr += "Development mode: Off\n"
	}

	if short {
		return versionStr
	}
	versionStr += version.Copyright + "\n"
	versionStr += version.GitHub + "\n\n"
	return versionStr
}
