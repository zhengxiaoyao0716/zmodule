package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/zhengxiaoyao0716/util/config"
	"github.com/zhengxiaoyao0716/zmodule/event"
	"github.com/zhengxiaoyao0716/zmodule/file"
	"github.com/zhengxiaoyao0716/zmodule/info"
)

// C .
type C map[string]interface{}

// GetString .
func (c *C) GetString(name string) string {
	if value, ok := (*cfg)[name]; ok {
		return value.(string)
	}
	return ""
}

// GetNum .
func (c *C) GetNum(name string) float64 {
	if value, ok := (*cfg)[name]; ok {
		return value.(float64)
	}
	return 0
}

// GetInt .
func (c *C) GetInt(name string) int {
	return int(c.GetNum(name))
}

// GetBool .
func (c *C) GetBool(name string) bool {
	if value, ok := (*cfg)[name]; ok {
		return value.(bool)
	}
	return false
}

var cfgPath string

// Dump .
func (c *C) Dump() error { return config.Dump(c, cfgPath) }

// Load .
func (c *C) Load(p string) error { return config.Load(p, c) }

var cfg *C

// Config return the instance of the inner config.
func Config() *C { return cfg }

// Load saved configure to the default config.
func Load() {
	cfg = config.LoadQ(cfgPath, &C{}).(*C)
	if cfg == nil {
		log.Fatalln("no inner launched config found.")
	}
	file.SetWorkDir(cfg.GetString("work_dir"))
}

// HasSavedCfg .
func HasSavedCfg() bool {
	if _, err := os.Lstat(cfgPath); err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// IsLog .
// func IsLog() bool { return cfg.Log != "" }
func IsLog() bool { return cfg.GetString("log") != "" }

// LogPath .
func LogPath() string {
	return file.AbsPath(cfg.GetString("log"), fmt.Sprintf(".%s.log", info.Name()))
}

func init() {
	event.OnInit(func(e event.Event) error {
		data := e.Data.(map[string]string)
		cfgPath = filepath.Join(data["workDir"], fmt.Sprintf(".%s.json", data["name"]))
		return nil
	})
}
