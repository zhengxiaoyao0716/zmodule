package main

import (
	"log"

	"github.com/kardianos/service"
	"github.com/zhengxiaoyao0716/zmodule"
)

// In this way that override those values,
// you can use `main` as the module name, instead of `github.com/zhengxiaoyao0716/zmodule`.
var (
	Version   string // `git describe --tags`
	Built     string // `date +%FT%T%z`
	GitCommit string // `git rev-parse --short HEAD`
	GoVersion string // `go version`
)

func main() {
	zmodule.Author = "zhengxiaoyao0716"
	zmodule.Homepage = "https://zhengxiaoyao0716.github.io/zmodule"
	zmodule.Repository = "https://github.com/zhengxiaoyao0716/zmodule"
	zmodule.License = "${Repository}/blob/master/LICENSE"

	zmodule.Version = Version
	zmodule.Built = Built
	zmodule.GitCommit = GitCommit
	zmodule.GoVersion = GoVersion

	zmodule.Main("zexample",
		&service.Config{
			Name:        "ZhengExampleService",
			DisplayName: "Example Service",
			Description: "Daemon service for example.",
		}, func() {
			log.Println("Running.")
		})
}
