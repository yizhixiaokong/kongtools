package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level        string `mapstructure:"level" yaml:"level"`
	AddSource    bool   `mapstructure:"addSource" yaml:"addSource"`
	Filename     string `mapstructure:"filename" yaml:"filename"`
	MaxSize      int    `mapstructure:"maxSize" yaml:"maxSize"`
	MaxBackups   int    `mapstructure:"maxBackups" yaml:"maxBackups"`
	MaxAge       int    `mapstructure:"maxAge" yaml:"maxAge"`
	RotateAtInit bool   `mapstructure:"rotateAtInit" yaml:"rotateAtInit"`
	MultiWriter  bool   `mapstructure:"multiWriter" yaml:"multiWriter"`
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		Level:        "debug",
		AddSource:    true,
		Filename:     "logs/kongtools.log",
		MaxSize:      10,
		MaxBackups:   3,
		MaxAge:       7,
		RotateAtInit: true,
		MultiWriter:  false,
	}
}

// InitLogger 初始化日志
func InitLogger(cfg Config) {
	var w io.Writer
	ll := &lumberjack.Logger{
		Filename:   cfg.Filename,   // filename
		MaxSize:    cfg.MaxSize,    // megabytes
		MaxAge:     cfg.MaxAge,     // days
		MaxBackups: cfg.MaxBackups, // max backups
	}

	if cfg.RotateAtInit {
		cobra.CheckErr(ll.Rotate()) // 启动时归档之前的日志
	}

	w = ll
	if cfg.MultiWriter {
		w = io.MultiWriter(os.Stdout, ll) // 同时写文件和屏幕 // !基本不需要
	}

	level := slog.LevelDebug
	cobra.CheckErr(level.UnmarshalText([]byte(cfg.Level)))

	logHandler := slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	})
	slog.SetDefault(slog.New(logHandler))

	// test log
	slog.Debug("debug log")
	slog.Info("info log")
	slog.Warn("warn log")
	slog.Error("error log")
}
