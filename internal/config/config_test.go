package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestConfig(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	
	once = sync.Once{}
	CfgFile = configPath

	cfg, err := Config()
	if err != nil {
		t.Fatalf("Config() failed: %v", err)
	}

	if cfg.App.TasksSavePath == "" {
		t.Error("TasksSavePath should not be empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should be created")
	}
}

func TestConfigWithInvalidPath(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test when running as root")
	}

	once = sync.Once{}
	CfgFile = "/root/no-permission-config/config.yaml"

	_, err := Config()
	if err == nil {
		t.Error("Should return error for inaccessible path")
	}
}
