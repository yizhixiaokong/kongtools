package paths

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
)

const appName = "kongtools"

var (
	// 可选覆盖路径（用于测试或自定义配置）
	overrideConfig string
	overrideData   string
	overrideCache  string
)

// Dirs 包含程序使用的三个基础路径。
type Dirs struct {
	Config string
	Data   string
	Cache  string
}

// All 返回当前使用的全部目录路径。
func All() Dirs {
	return Dirs{
		Config: ConfigDir(),
		Data:   DataDir(),
		Cache:  CacheDir(),
	}
}

// Override 允许在运行时（如测试）覆盖默认路径。
// 传入空字符串则忽略。
func Override(config, data, cache string) {
	if config != "" {
		overrideConfig = config
	}
	if data != "" {
		overrideData = data
	}
	if cache != "" {
		overrideCache = cache
	}
}

// --- 路径获取函数 ---

func ConfigDir() string {
	return ensureDir(resolvePath(
		overrideConfig,
		os.Getenv("KONGTOOLS_CONFIG_HOME"),
		filepath.Join(xdg.ConfigHome, appName),
	))
}

func DataDir() string {
	return ensureDir(resolvePath(
		overrideData,
		os.Getenv("KONGTOOLS_DATA_HOME"),
		filepath.Join(xdg.DataHome, appName),
	))
}

func CacheDir() string {
	return ensureDir(resolvePath(
		overrideCache,
		os.Getenv("KONGTOOLS_CACHE_HOME"),
		filepath.Join(xdg.CacheHome, appName),
	))
}

// --- 文件路径辅助 ---

func ConfigFile(name string) string {
	return filepath.Join(ConfigDir(), name)
}

func DataFile(name string) string {
	return filepath.Join(DataDir(), name)
}

func CacheFile(name string) string {
	return filepath.Join(CacheDir(), name)
}

// --- 实用函数 ---

// ClearCache 清空缓存目录内容（不会删除目录本身）。
func ClearCache() error {
	cache := CacheDir()
	entries, err := os.ReadDir(cache)
	if err != nil {
		return err
	}
	for _, e := range entries {
		os.RemoveAll(filepath.Join(cache, e.Name()))
	}
	return nil
}

// resolvePath 返回优先可用的路径（从 override → env → default）。
func resolvePath(paths ...string) string {
	for _, p := range paths {
		if p != "" {
			return p
		}
	}
	return ""
}

// ensureDir 确保目录存在。
func ensureDir(dir string) string {
	_ = os.MkdirAll(dir, 0o755) // 幂等操作，多次调用安全
	return dir
}
