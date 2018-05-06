package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/zhengxiaoyao0716/zmodule/event"
)

// MoveAway can to move away the file instead of delete it, used to avoid the file exist conflict.
func MoveAway(fp string) error {
	if _, err := os.Lstat(fp); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		if !os.IsExist(err) {
			return err
		}
	}
	ext := ""
	for i := len(fp) - 1; i >= 0 && !os.IsPathSeparator(fp[i]); i-- {
		if fp[i] == '.' {
			ext = fp[i:]
			fp = fp[0:i]
			break
		}
	}
	err := os.Rename(fp+ext, fmt.Sprintf("%s.%d%s", fp, time.Now().Unix(), ext))
	if err != nil {
		return err
	}
	return nil
}

var (
	workDir string
)

// WorkDir return the really directory path where use work in.
// Other ways such as os.Getwd, os.Executable doesn't work in daemon service.
func WorkDir() string { return workDir }

// SetWorkDir .
func SetWorkDir(dir string) { workDir = AbsPath(dir) }

// AbsPath is the absolute path build with WorkDir.
// If you notice that, you would know that it doesn't support unsafe up-find, such as "./../".
// If you need to use a directory or file out of current, you need to input absolute path.
func AbsPath(elem ...string) string {
	p := path.Join(elem...)
	if filepath.IsAbs(p) {
		return p
	}
	return path.Join(workDir, p)
}

func init() {
	event.OnInit(func(e event.Event) error {
		data := e.Data.(map[string]string)
		workDir = data["workDir"]
		return nil
	})
}
