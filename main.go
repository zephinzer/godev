//go:generate go run data/generate.go
package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"

	shellquote "github.com/kballard/go-shellquote"
)

func main() {
	config := InitConfig()
	godev := InitGoDev(config)
	godev.Start()
}

type InitFile struct {
	Filename string
	Data     string
}

var InitFileMap = map[string]*InitFile{
	"dockerfile":    &InitFile{"Dockerfile", DataDockerfile},
	"makefile":      &InitFile{"Makefile", DataMakefile},
	".dockerignore": &InitFile{".dockerignore", DataDotDockerignore},
	".gitignore":    &InitFile{".gitignore", DataDotGitignore},
	"main.go":       &InitFile{"main.go", DataMainDotgo},
	"go.mod":        &InitFile{"go.mod", DataGoDotMod},
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
	config  *Config
	logger  *Logger
	watcher *Watcher
	runner  *Runner
}

func (godev *GoDev) Start() {
	defer godev.logger.Infof("godev has ended")
	godev.logger.Infof("godev has started")
	if godev.config.RunVersion {
		fmt.Printf("%s-%s\n", Version, Commit)
	} else if godev.config.RunView {
		godev.viewFile()
	} else if godev.config.RunInit {
		godev.initialiseDirectory()
	} else {
		godev.startWatching()
	}
}

func (godev *GoDev) createPipeline() []*ExecutionGroup {
	var pipeline []*ExecutionGroup
	for _, execGroup := range godev.config.ExecGroups {
		executionGroup := &ExecutionGroup{}
		var executionCommands []*Command
		commands := strings.Split(execGroup, godev.config.CommandsDelimiter)
		for _, command := range commands {
			if sections, err := shellquote.Split(command); err != nil {
				panic(err)
			} else {
				executionCommands = append(
					executionCommands,
					InitCommand(&CommandConfig{
						Application: sections[0],
						Arguments:   sections[1:],
						Directory:   godev.config.WorkDirectory,
						LogLevel:    godev.config.LogLevel,
					}),
				)
			}
		}
		executionGroup.commands = executionCommands
		pipeline = append(pipeline, executionGroup)
	}
	return pipeline
}

func (godev *GoDev) eventHandler(events *[]WatcherEvent) bool {
	for _, e := range *events {
		godev.logger.Trace(e)
	}
	godev.runner.Trigger()
	return true
}

func (godev *GoDev) startWatching() {
	godev.logUniversalConfigurations()
	godev.logWatchModeConfigurations()
	godev.initialiseWatcher()
	godev.initialiseRunner()

	var wg sync.WaitGroup
	godev.watcher.BeginWatch(&wg, godev.eventHandler)
	godev.logger.Infof("working dir : '%s'", godev.config.WorkDirectory)
	godev.logger.Infof("watching dir: '%s'", godev.config.WatchDirectory)
	godev.runner.Trigger()
	wg.Wait()
}

func (godev *GoDev) logUniversalConfigurations() {
	godev.logger.Debugf("flag - init       : %v", godev.config.RunInit)
	godev.logger.Debugf("flag - test       : %v", godev.config.RunTest)
	godev.logger.Debugf("flag - view       : %v", godev.config.RunView)
	godev.logger.Debugf("watch directory   : %s", godev.config.WatchDirectory)
	godev.logger.Debugf("work directory    : %s", godev.config.WorkDirectory)
	godev.logger.Debugf("build output      : %s", godev.config.BuildOutput)
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

func (godev *GoDev) initialiseInitialisers() []Initialiser {
	return []Initialiser{
		InitGitInitialiser(&GitInitialiserConfig{
			Path: path.Join(godev.config.WorkDirectory),
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/.gitignore"),
			Data:     []byte(DataDotGitignore),
			Question: "seed a .gitignore?",
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/go.mod"),
			Data:     []byte(DataGoDotMod),
			Question: "seed a go.mod?",
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/main.go"),
			Data:     []byte(DataMainDotgo),
			Question: "seed a main.go?",
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/Dockerfile"),
			Data:     []byte(DataDockerfile),
			Question: "seed a Dockerfile?",
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/.dockerignore"),
			Data:     []byte(DataDotDockerignore),
			Question: "seed a .dockerignore?",
		}),
		InitFileInitialiser(&FileInitialiserConfig{
			Path:     path.Join(godev.config.WorkDirectory, "/Makefile"),
			Data:     []byte(DataMakefile),
			Question: "seed a Makefile?",
		}),
	}
}

// initialiseDirectory assists in initialising the working directory
func (godev *GoDev) initialiseDirectory() {
	if !directoryExists(godev.config.WorkDirectory) {
		godev.logger.Errorf("the directory at '%s' does not exist - create it first with:\n  mkdir -p %s", godev.config.WorkDirectory, godev.config.WorkDirectory)
		os.Exit(1)
	}
	initialisers := godev.initialiseInitialisers()
	for i := 0; i < len(initialisers); i++ {
		initialiser := initialisers[i]
		if initialiser.Check() {
			err := initialiser.Handle(true)
			if err != nil {
				fmt.Println(Color("red", err.Error()))
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			if initialiser.Confirm(reader) {
				fmt.Println(Color("green", "godev> sure thing"))
				initialiser.Handle()
			} else {
				fmt.Println(Color("yellow", "godev> lets skip that then"))
			}
		}
	}
}

func (godev *GoDev) initialiseRunner() {
	godev.runner = InitRunner(&RunnerConfig{
		Pipeline: godev.createPipeline(),
		LogLevel: godev.config.LogLevel,
	})
}

func (godev *GoDev) initialiseWatcher() {
	godev.watcher = InitWatcher(&WatcherConfig{
		FileExtensions: godev.config.FileExtensions,
		IgnoredNames:   godev.config.IgnoredNames,
		RefreshRate:    godev.config.Rate,
		LogLevel:       godev.config.LogLevel,
	})
	godev.watcher.RecursivelyWatch(godev.config.WatchDirectory)
}

func (godev *GoDev) viewFile() {
	config := godev.config
	logger := godev.logger
	fileKey := strings.ToLower(config.View)
	if InitFileMap[fileKey] != nil {
		initFile := InitFileMap[fileKey]
		logger.Infof("previewing contents of %s", initFile.Filename)
		fmt.Println(initFile.Filename)
		logger.Infof("end of preview for contents of %s", initFile.Filename)
	} else {
		logger.Panicf("the requested file '%s' does not seem to exist :/", fileKey)
	}
}
