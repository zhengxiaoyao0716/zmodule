// Package info helper provide some global static single instance.
package info

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zhengxiaoyao0716/util/cout"

	"github.com/zhengxiaoyao0716/zmodule/event"
)

type info struct {
	Name    string
	WorkDir string

	Version string

	Author     string
	Homepage   string
	Repository string
	License    string

	Built     string
	GitCommit string
	GoVersion string
}

var i info

// Name .
func Name() string { return i.Name }

// WorkDir .
func WorkDir() string { return i.WorkDir }

// Info .
func Info() string {
	var (
		kl, vl string
		kli    = 6
		vli    = 12
		lines  = []func() string{
			func() string {
				return "--" + strings.Repeat("-", kli) + "---" + strings.Repeat("-", vli) + "--"
			},
		}
		max = func(l, r int) int {
			if l > r {
				return l
			}
			return r
		}
	)
	lines = append(lines, func() string {
		li := kli + vli + 3
		l := strconv.Itoa(li)

		nv := i.Name + " " + i.Version
		a := i.Author

		d := li - (len(nv) + len(a))
		if d > 0 {
			if d > 8 {
				d -= 8
				a = "author: " + a
			}
			return "| " + cout.Info(nv+strings.Repeat(" ", d)+a) + " |"
		}
		if li-len(a) > 8 {
			a = "author: " + a
		}
		return "| " + cout.Info(fmt.Sprintf("%-"+l+"s", nv)) + " |\n" +
			"| " + cout.Info(fmt.Sprintf("%+"+l+"s", a)) + " |"
	})
	lines = append(lines, lines[0])
	for _, kv := range [][2]string{
		{"Homepage", i.Homepage},
		{"Repository", i.Repository},
		{"License", i.License},
		{"", ""},
		{"Built", i.Built},
		{"GitCommit", i.GitCommit},
		{"GoVersion", i.GoVersion},
	} {
		func(k, v string) {
			if k == "" {
				lines = append(lines, func() string {
					return "| " + strings.Repeat("-", kli) + " | " + strings.Repeat("-", vli) + " |"
				})
				return
			}
			kli = max(kli, len(k))
			vli = max(vli, len(v))
			lines = append(lines, func() string {
				return "| " + cout.Info(fmt.Sprintf("%+"+kl+"s", k)) + " | " + cout.Info(fmt.Sprintf("%-"+vl+"s", v)) + " |"
			})
		}(kv[0], kv[1])
	}
	vli = max(kli+kli>>1, vli)
	kl = strconv.Itoa(kli)
	vl = strconv.Itoa(vli)
	lines = append(lines, lines[0])
	return fmt.Sprintln(
		strings.Join(func() []string {
			var result []string
			for _, line := range lines {
				result = append(result, line())
			}
			return result
		}(), "\n"),
	)
}

// Version .
func Version() string { return fmt.Sprintf("Version: %s\n", cout.Info(i.Version)) }

func init() {
	event.OnInit(func(e event.Event) error {
		data := e.Data.(map[string]string)

		i.Name = data["name"]
		i.WorkDir = data["workDir"]

		i.Version = data["Version"]
		i.Author = data["Author"]
		i.Homepage = data["Homepage"]
		i.Repository = data["Repository"]
		i.License = data["License"]
		i.Built = data["Built"]
		i.GitCommit = data["GitCommit"]
		i.GoVersion = data["GoVersion"]

		return nil
	})
}
