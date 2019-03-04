//go:generate go run data/generate.go
package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	shellquote "github.com/kballard/go-shellquote"
)

// Version should be populated with -ldflags on build with the semver version
var Version string

// Commit should be populated with -ldflags on build with the current git commit
var Commit string

func main() {
	config := InitConfig()
	godev := InitGoDev(config)
	godev.Start()
}

func InitGoDev(config *Config) *GoDev {
	return &GoDev{
		config: config,
		logger: InitLogger(&LoggerConfig{
			Name:   "main",
			Format: "production",
			Level:  config.LogLevel,
		}),
	}
}

type GoDev struct {
	config *Config
	logger *Logger
}

func (godev *GoDev) Start() {
	defer godev.logger.Infof("godev has ended")
	godev.logger.Infof("godev has started")
	if godev.config.RunVersion {
		fmt.Printf("godev %s-%s\n", Version, Commit)
	} else if godev.config.RunView {
		godev.viewFile()
	} else if godev.config.RunInit {
		godev.initialiseDirectory()

	} else {
		godev.startWatching()
	}
}

func (godev *GoDev) startWatching() {
	config := godev.config
	logger := godev.logger
	godev.logUniversalConfigurations()
	godev.logWatchModeConfigurations()
	watcher := InitWatcher(&WatcherConfig{
		FileExtensions: config.FileExtensions,
		IgnoredNames:   config.IgnoredNames,
		RefreshRate:    config.Rate,
		LogLevel:       config.LogLevel,
	})
	watcher.RecursivelyWatch(config.WatchDirectory)
	var pipeline []*ExecutionGroup
	for _, execGroup := range config.ExecGroups {
		executionGroup := &ExecutionGroup{}
		var executionCommands []*Command
		commands := strings.Split(execGroup, config.CommandsDelimiter)
		for _, command := range commands {
			if sections, err := shellquote.Split(command); err != nil {
				panic(err)
			} else {
				executionCommands = append(
					executionCommands,
					InitCommand(&CommandConfig{
						Application: sections[0],
						Arguments:   sections[1:],
						Directory:   config.WatchDirectory,
						LogLevel:    config.LogLevel,
					}),
				)
			}
		}
		executionGroup.commands = executionCommands
		pipeline = append(pipeline, executionGroup)
	}
	runner := InitRunner(&RunnerConfig{
		Pipeline: pipeline,
		LogLevel: config.LogLevel,
	})

	var wg sync.WaitGroup
	watcher.BeginWatch(&wg, func(events *[]WatcherEvent) bool {
		for _, e := range *events {
			logger.Trace(e)
		}
		runner.Trigger()
		return true
	})

	logger.Infof("started watcher at %s", config.WatchDirectory)

	runner.Trigger()
	wg.Wait()
}

func (godev *GoDev) logUniversalConfigurations() {
	config := godev.config
	logger := godev.logger
	logger.Debugf("flag - init       : %v", config.RunInit)
	logger.Debugf("flag - test       : %v", config.RunTest)
	logger.Debugf("flag - view       : %v", config.RunView)
	logger.Debugf("watch directory   : %s", config.WatchDirectory)
	logger.Debugf("build output      : %s", config.BuildOutput)
}

func (godev *GoDev) logWatchModeConfigurations() {
	config := godev.config
	logger := godev.logger
	logger.Debugf("file extensions   : %v", config.FileExtensions)
	logger.Debugf("ignored names     : %v", config.IgnoredNames)
	logger.Debugf("refresh interval  : %v", config.Rate)
	logger.Debugf("execution delim   : %s", config.CommandsDelimiter)
	logger.Debug("execution groups as follows...")
	for egIndex, execGroup := range config.ExecGroups {
		logger.Debugf("  %v) %s", egIndex+1, execGroup)
		commands := strings.Split(execGroup, config.CommandsDelimiter)
		for cIndex, command := range commands {
			sections, err := shellquote.Split(command)
			if err != nil {
				panic(err)
			}
			app := sections[0]
			args := sections[1:]
			logger.Debugf("    %v > %s %v", cIndex+1, app, args)
		}
	}
}

func (godev *GoDev) initialiseDirectory() {
	types := []string{"git", ".gitignore", "go.mod", "main.go", "Dockerfile", ".dockerignore", "Makefile"}
	checks := map[string]func() bool{
		"git": func() bool {
			return directoryExists(path.Join(godev.config.WatchDirectory, "/.git"))
		},
		".gitignore": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/.gitignore"))
		},
		"go.mod": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/go.mod"))
		},
		"main.go": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/main.go"))
		},
		"Dockerfile": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/Dockerfile"))
		},
		".dockerignore": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/.dockerignore"))
		},
		"Makefile": func() bool {
			return fileExists(path.Join(godev.config.WatchDirectory, "/Makefile"))
		},
	}
	questions := map[string]string{
		"git":           "initialise git repository?",
		".gitignore":    "seed a .gitignore?",
		"go.mod":        "seed a go.mod?",
		"main.go":       "seed a main.go?",
		"Dockerfile":    "seed a Dockerfile?",
		".dockerignore": "seed a .dockerignore?",
		"Makefile":      "seed a Makefile?",
	}
	handlers := map[string]func(...bool) error{
		"git": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping initialisation of git repository - already initialised"))
				return nil
			}
			fmt.Println("git todo")
			return errors.New("todo")
		},
		".gitignore": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of .gitignore - already exists"))
				return nil
			}
			fmt.Println(".gitignore todo")
			return errors.New("todo")
		},
		"go.mod": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of go.mod - already exists"))
				return nil
			}
			fmt.Println("go.mod todo")
			return errors.New("todo")
		},
		"main.go": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of main.go - already exists"))
				return nil
			}
			fmt.Println("main.go todo")
			return errors.New("todo")
		},
		"Dockerfile": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of Dockerfile - already exists"))
				return nil
			}
			fmt.Println("Dockerfile todo")
			return errors.New("todo")
		},
		".dockerignore": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of .dockerignore - already exists"))
				return nil
			}
			filePath := path.Join(getCurrentWorkingDirectory(), "/.dockerignore")
			file, err := os.Create(filePath)
			if err != nil {
				return err
			}
			size, err := file.Write([]byte(DataDotDockerignore))
			if err != nil {
				return err
			}
			fmt.Println(Color("green", fmt.Sprintf("godev> written %v bytes to %s", size, filePath)))
			return errors.New("todo")
		},
		"Makefile": func(skip ...bool) error {
			if len(skip) > 0 && skip[0] {
				fmt.Println(Color("gray", "godev> skipping seeding of Makefile - already exists"))
				return nil
			}
			fmt.Println("Makefile todo")
			return errors.New("todo")
		},
	}
	for i := 0; i < len(types); i++ {
		id := types[i]
		if checks[id]() {
			err := handlers[id](true)
			if err != nil {
				fmt.Println(Color("red", err.Error()))
			}
		} else {
			c := confirm(
				Color("white", "godev> "+questions[id]),
				false,
				Color("bold", Color("red", "sorry, i didn't get that")),
			)
			if *c {
				fmt.Println(Color("green", "godev> sure thing"))
				handlers[id]()
			} else {
				fmt.Println(Color("yellow", "godev> lets skip that then"))
			}
		}
	}
}

func (godev *GoDev) viewFile() {
	config := godev.config
	logger := godev.logger
	switch strings.ToLower(config.View) {
	case "dockerfile":
		logger.Info("previewing contents of Dockerfile")
		fmt.Println(DataDockerfile)
		logger.Info("end of preview for contents of Dockerfile")
	case "makefile":
		logger.Info("previewing contents of Makefile")
		fmt.Println(DataMakefile)
		logger.Info("end of preview for contents of Makefile")
	case ".dockerignore":
		logger.Info("previewing contents of .dockerignore")
		fmt.Println(DataDotDockerignore)
		logger.Info("end of preview for contents of .dockerignore")
	case ".gitignore":
		logger.Info("previewing contents of .gitignore")
		fmt.Println(DataDotGitignore)
		logger.Info("end of preview for contents of .gitignore")
	case "main.go":
		logger.Info("previewing contents of main.go")
		fmt.Println(DataMainDotgo)
		logger.Info("end of preview for contents of main.go")
	case "go.mod":
		logger.Info("previewing contents of go.mod")
		fmt.Println(DataGoDotMod)
		logger.Info("end of preview for contents of go.mod")
	}
}
