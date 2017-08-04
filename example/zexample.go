package main

import (
	"log"

	"github.com/kardianos/service"
	"github.com/zhengxiaoyao0716/zmodule"
)

func main() {
	zmodule.Author = "zhengxiaoyao0716"
	zmodule.Homepage = "https://zhengxiaoyao0716.github.io/zmodule"
	zmodule.Repository = "https://github.com/zhengxiaoyao0716/zmodule"
	zmodule.License = "${Repository}/blob/master/LICENSE"

	zmodule.Main("zexample",
		&service.Config{
			Name:        "ZhengExampleService",
			DisplayName: "Example Service",
			Description: "Daemon service for example.",
		}, func() {
			log.Println("Running.")
		})
}
