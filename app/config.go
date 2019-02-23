package main

import (
	"flag"
)

type Config struct {
	RunInit     bool
	RunTest     bool
	RunModWatch bool
}

func InitConfig() *Config {
	config := &Config{}
	flag.BoolVar(&config.RunInit, "init", false, "when this flag is specified, godev initiaises the current directory")
	flag.BoolVar(&config.RunTest, "test", false, "when this flag is specified, godev runs the tests with coverage")
	flag.BoolVar(&config.RunModWatch, "watch", false, "when this flag is specified, godev runs the command in watch mode if applicable")
	flag.Parse()
	return config
}
