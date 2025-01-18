package internal

import (
	"fmt"
	"runtime"
)

var (
	buildVersion = "dev"
	buildCommit  = "n/a"
)

type VersionInfo struct {
	BuildVersion string `json:"buildVersion" yaml:"buildVersion"`
	BuildCommit  string `json:"buildCommit" yaml:"buildCommit"`
	Platform     string `json:"platform" yaml:"platform"`
}

func GetVersion() *VersionInfo {
	return &VersionInfo{
		BuildVersion: buildVersion,
		BuildCommit:  buildCommit,
		Platform:     fmt.Sprintf("%s/%s, %s", runtime.GOOS, runtime.GOARCH, runtime.Version()),
	}
}
