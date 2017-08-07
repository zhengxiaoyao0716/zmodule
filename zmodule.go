package zmodule

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kardianos/service"

	"github.com/zhengxiaoyao0716/util/console"
	"github.com/zhengxiaoyao0716/util/flag"
	"github.com/zhengxiaoyao0716/zmodule/config"
	"github.com/zhengxiaoyao0716/zmodule/event"
	"github.com/zhengxiaoyao0716/zmodule/file"
	"github.com/zhengxiaoyao0716/zmodule/info"
)

type program struct {
	run func()
}

func (p *program) Start(s service.Service) error {
	// Load saved config.
	config.Load()

	// Redirect logs.
	if config.IsLog() || !service.Interactive() {
		// file, err := os.OpenFile(config.LogPath(), os.O_APPEND|os.O_CREATE, 0666)
		fp := config.LogPath()

		if err := file.MoveAway(fp); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(fp), 0600); err != nil {
			return err
		}
		file, err := os.Create(fp)
		if err != nil {
			return err
		}
		// close log file before stop.
		event.On(event.KeyStop, func(event.Event) error { return file.Close() })

		log.SetOutput(file)
	}
	event.Emit(event.KeyStart, config.Config())
	event.Pool().Wait()
	// Start should not block. Do the actual work async.
	log.Println("Service start.")
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	log.Println("Service stop.")
	event.Emit(event.KeyStop, nil)
	event.Pool().Wait()
	return nil
}

// Argument .
type Argument struct {
	Type    string // string, int, bool
	Default interface{}
	Usage   string
}

// Args .
// You can add your custom args directly, and then they would be parsed into config.
// Please don't change both the name and the type of the args reserved below.
var Args = map[string]Argument{
	"config":   {"string", "", "Config json file path."},
	"work_dir": {"string", "", "Directory to find or storage files."},
	"log":      {"string", "", "Path to storage logger files."},
}

// parseFlag parsed the remained args, load them to config, then dump the config to save them.
func parseFlag(args []string) {
	flags := map[string]func() interface{}{}
	for name, arg := range Args {
		switch arg.Type {
		case "string":
			f := flag.String(name, arg.Usage)
			flags[name] = func() interface{} {
				if *f == nil {
					return nil
				}
				return (*f)()
			}
		case "int":
			f := flag.Int(name, arg.Usage)
			flags[name] = func() interface{} {
				if *f == nil {
					return nil
				}
				return (*f)()
			}
		case "bool":
			f := flag.Bool(name, arg.Usage)
			flags[name] = func() interface{} {
				if *f == nil {
					return nil
				}
				return (*f)()
			}
		}
	}

	flag.CommandLine.Parse(args)

	cfgPath := flags["config"]()
	cfg := &config.C{}
	if cfgPath != nil {
		if err := cfg.Load(cfgPath.(string)); err != nil {
			log.Fatalln(err)
		}
	}

	for name, fn := range flags {
		// Launch arguments
		if value := fn(); value != nil {
			(*cfg)[name] = value
			continue
		}
		// User config
		if _, ok := (*cfg)[name]; ok {
			continue
		}
		// Os environment
		if value := os.Getenv(info.Name() + "_" + name); value != "" {
			(*cfg)[name] = value
			continue
		}
		// Default value in code
		if value := Args[name].Default; value != nil {
			(*cfg)[name] = value
			continue
		}
		// Not found
		log.Fatalln(errors.New("missing runner argument: name=" + name))
	}

	if err := cfg.Dump(); err != nil {
		log.Fatalln(err)
	}
}

// Command .
type Command struct {
	Usage   string
	Handler func(string, []string)
}

var srv service.Service

// Cmds .
var Cmds = map[string]Command{
	"(NIL)": Command{
		"Quick run with last args.",
		func(parsed string, args []string) {
			if !config.HasSavedCfg() {
				parseFlag(args)
			}
			if err := srv.Run(); err != nil {
				// Redirect log out to a temp file and write the error.
				if !service.Interactive() {
					file, err := os.OpenFile(filepath.Join(info.WorkDir(), "error.log"), os.O_APPEND|os.O_CREATE, 0666)
					if err != nil {
						// 这尼玛如果都出错了我真tm不知道怎么报告给用户了
					}
					defer file.Close()

					log.SetOutput(file)
				}

				log.Fatalln(err)
			}
		}},
	"run": Command{
		"Run.",
		func(parsed string, args []string) {
			parseFlag(args)
			if err := srv.Run(); err != nil {
				log.Fatalln(err)
			}
		}},
	"version": Command{
		"Show the version.",
		func(string, []string) { fmt.Println(info.Version()) }},
	"info": Command{
		"Show the info.",
		func(string, []string) { fmt.Println(info.Info()) }},
	"service": Command{
		"Control the system service.",
		func(parsed string, args []string) {
			control := func(arg string) {
				err := service.Control(srv, arg)
				if err != nil {
					log.Fatal(err)
				}
			}

			if len(args) == 0 {
				args = []string{"(NIL)"}
			}

			ParseCmd(parsed, args, map[string]Command{
				"(NIL)": Command{
					"Quick start with last args.",
					func(parsed string, args []string) {
						if !config.HasSavedCfg() {
							parseFlag(args)
						}
						control("start")
					}},
				"start": Command{
					"Start the service.",
					func(parsed string, args []string) {
						parseFlag(args)
						control("start")
					}},
				"stop": Command{
					"Stop the service.",
					func(string, []string) { control("stop") }},
				"restart": Command{
					"Restart the service.",
					func(string, []string) { control("restart") }},
				"install": Command{
					"Install the service.",
					func(string, []string) { control("install") }},
				"uninstall": Command{
					"Uninstall the service.",
					func(string, []string) { control("uninstall") }},
			})
		}},
}

// ParseCmd .
func ParseCmd(parsed string, args []string, cmds map[string]Command) {
	if len(args) > 0 {
		console.PushLine(args[0])
		args = args[1:]
	}

	linked := []string{}
	for arg := range cmds {
		linked = append(linked, arg)
	}

	sort.Strings(linked)
	arg := strings.TrimSpace(console.ReadWord(fmt.Sprintf("Usage: %s\n(You can also enter '--help' to check details)\n> ", strings.Join(linked, " "))))
	switch arg {
	default:
		if cmd, ok := cmds[arg]; ok {
			cmd.Handler(parsed+" "+arg, args)
			return
		}
		fmt.Printf("%s: invalid option: '%s' for Command '%s'.\n", info.Name(), arg, parsed)
		fallthrough
	case "--help":
		fallthrough
	case "-h":
		fmt.Println("Usage:", parsed, "<Command>")
		fmt.Println()
		fmt.Println("Commands list:")
		for _, arg := range linked {
			cmd := cmds[arg]
			fmt.Printf("    %10s\t%s\n", arg, cmd.Usage)
		}
	}
}

// Those values should be assigned during compile statement.
// ``` bash
// m=github.com/zhengxiaoyao0716/zmodule
// v=`git describe --tags`
// b=`date +%FT%T%z`
// go run -ldflags "-X $m.Version=$v -X $m.Built=$b ..." zexample.go version
// ```
// But some of them may not need to changed frequently.
// So you can also set them in your code directly.
var (
	Version string // `git describe --tags`

	Author     string // zhengxiaoyao0716
	Homepage   string // https://zhengxiaoyao0716.github.io/zmodule
	Repository string // https://github.com/zhengxiaoyao0716/zmodule
	License    string // https://github.com/zhengxiaoyao0716/zmodule/blob/master/LICENSE

	Built     string // `date +%FT%T%z`
	GitCommit string // `git rev-parse --short HEAD`
	GoVersion string // `go version`
)

// Main .
func Main(name string, scfg *service.Config, run func(), cusKs ...[2]string) {
	event.Init(map[string]string{
		"name": name,
		"workDir": func() string {
			var (
				workDir string
				err     error
			)
			if service.Interactive() {
				workDir, err = os.Getwd()
			} else {
				workDir, err = filepath.Abs(filepath.Dir(os.Args[0]))
			}
			if err != nil {
				log.Fatal(err)
			}
			return workDir
		}(),
		"Version":    Version,
		"Author":     Author,
		"Homepage":   Homepage,
		"Repository": Repository,
		"License":    License,
		"Built":      Built,
		"GitCommit":  GitCommit,
		"GoVersion":  GoVersion,
	}, cusKs...)
	event.Pool().Wait()

	var err error
	srv, err = service.New(&program{run}, scfg)
	if err != nil {
		log.Fatal(err)
	}

	args := os.Args[1:]
	if len(args) == 0 {
		args = []string{"(NIL)"}
	}
	ParseCmd(name, args, Cmds)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
