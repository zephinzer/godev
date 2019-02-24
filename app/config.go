package main

import (
	"flag"
	"time"
)

const DefaultRefreshRate = 2 * time.Second

type Config struct {
	RunInit        bool
	RunTest        bool
	RunModWatch    bool
	WatchDirectory string
	Rate           time.Duration
}

func InitConfig() *Config {
	currentWorkingDirectory := getCurrentWorkingDirectory()
	config := &Config{}
	flag.StringVar(&config.WatchDirectory, "dir", currentWorkingDirectory, "specifies the directory to watch")
	flag.DurationVar(&config.Rate, "rate", DefaultRefreshRate, "specifies the refresh rate of the file system watch")
	flag.BoolVar(&config.RunInit, "init", false, "when this flag is specified, godev initiaises the current directory")
	flag.BoolVar(&config.RunTest, "test", false, "when this flag is specified, godev runs the tests with coverage")
	flag.BoolVar(&config.RunModWatch, "watch", false, "when this flag is specified, godev runs the command in watch mode if applicable")
	flag.Parse()
	return config
}
