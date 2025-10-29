package tui

import "kongtools/internal/pkg/paths"

// Config TUI 配置
type Config struct {
	TasksSavePath string
}

// DefaultConfig 默认配置 YAML
const DefaultConfig = `app:
  tasksSavePath: 
`

// NewConfig 创建默认配置
func NewConfig() Config {
	return Config{
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
