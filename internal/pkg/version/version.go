package version

import (
	"fmt"
	"runtime"
)

var (
	// Version 版本号，可在编译时通过 -ldflags 设置
	Version = "dev"
	// Commit Git 提交 hash，可在编译时通过 -ldflags 设置
	Commit = "unknown"
	// BuildTime 编译时间，可在编译时通过 -ldflags 设置
	BuildTime = "unknown"
	// GoVersion Go 版本
	GoVersion = runtime.Version()
)

// GetVersion 返回版本信息
func GetVersion() string {
	return Version
}

// GetCommit 返回 Git 提交 hash
func GetCommit() string {
	return Commit
}

// GetBuildTime 返回编译时间
func GetBuildTime() string {
	return BuildTime
}

// GetGoVersion 返回 Go 版本
func GetGoVersion() string {
	return GoVersion
}

// GetFullVersion 返回完整的版本信息
func GetFullVersion() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuild Time: %s\nGo Version: %s",
		Version, Commit, BuildTime, GoVersion)
}

// GetVersionInfo 返回版本信息的结构化数据
func GetVersionInfo() map[string]string {
	return map[string]string{
		"version":   Version,
		"commit":    Commit,
		"buildTime": BuildTime,
		"goVersion": GoVersion,
	}
}
