package file

import (
	"path"
	"path/filepath"

	"github.com/zhengxiaoyao0716/zmodule/event"
)

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
