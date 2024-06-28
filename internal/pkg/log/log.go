package log

import (
	"io"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Level        string
	AddSource    bool
	Filename     string
	MaxSize      int
	MaxBackups   int
	MaxAge       int
	RotateAtInit bool
	MultiWriter  bool
}

const DefaultConfig = `log:
  level: debug
  addSource: true
  filename: kongtools.log
  maxSize: 10
  maxBackups: 3
  maxAge: 7
  rotateAtInit: true
  multiWriter: false
`

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
