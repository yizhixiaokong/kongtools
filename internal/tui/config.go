package tui

import "kongtools/internal/pkg/paths"

// Config TUI 配置
type Config struct {
	TasksSavePath string `mapstructure:"tasksSavePath" yaml:"tasksSavePath"`
}

// Default returns a Config with default values
func Default() *Config {
	return &Config{
		TasksSavePath: paths.DataFile("tasks.json"),
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.TasksSavePath == "" {
		c.TasksSavePath = paths.DataFile("tasks.json")
	}
	return nil
}
