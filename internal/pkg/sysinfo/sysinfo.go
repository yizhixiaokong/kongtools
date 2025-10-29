package sysinfo

import (
	"kongtools/internal/config"
	"kongtools/internal/pkg/paths"
	"kongtools/internal/pkg/version"
	"log/slog"
	"runtime"

	"github.com/spf13/viper"
)

// PrintSystemInfo 打印系统信息到日志
// 应该在所有初始化完成后调用
func PrintSystemInfo() {
	slog.Info("========== System Information ==========")

	// 打印版本信息
	slog.Info("Application version",
		slog.String("version", version.GetVersion()),
		slog.String("commit", version.GetCommit()),
		slog.String("buildTime", version.GetBuildTime()),
		slog.String("goVersion", version.GetGoVersion()),
	)

	// 打印配置文件路径
	slog.Info("Config file", slog.String("path", viper.ConfigFileUsed()))

	// 打印目录信息
	dirs := paths.All()
	// slog.Info("Config directory", slog.String("path", dirs.Config))
	slog.Info("Data directory", slog.String("path", dirs.Data))
	slog.Info("Cache directory", slog.String("path", dirs.Cache))

	// 打印日志配置信息
	cfg := config.Config()
	slog.Info("Log configuration",
		slog.String("level", cfg.Log.Level),
		slog.Bool("addSource", cfg.Log.AddSource),
		slog.String("filename", cfg.Log.Filename),
		slog.Int("maxSize", cfg.Log.MaxSize),
		slog.Int("maxBackups", cfg.Log.MaxBackups),
		slog.Int("maxAge", cfg.Log.MaxAge),
		slog.Bool("rotateAtInit", cfg.Log.RotateAtInit),
		slog.Bool("multiWriter", cfg.Log.MultiWriter),
	)

	// 打印 Go 运行时信息
	slog.Info("Go runtime",
		slog.String("version", runtime.Version()),
		slog.String("os", runtime.GOOS),
		slog.String("arch", runtime.GOARCH),
	)

	slog.Info("========================================")
}
