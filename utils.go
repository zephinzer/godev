package main

import (
	"os"
	"strings"
)

type ConfigCommaDelimitedString []string

func (ccds *ConfigCommaDelimitedString) Set(item string) error {
	*ccds = append(*ccds, strings.Split(item, ",")...)
	return nil
}

func (ccds *ConfigCommaDelimitedString) String() string {
	return strings.Join(*ccds, ",")
}

type ConfigMultiflagString []string

func (cmfs *ConfigMultiflagString) Set(item string) error {
	*cmfs = append(*cmfs, item)
	return nil
}

func (cmfs *ConfigMultiflagString) String() string {
	return strings.Join(*cmfs, ",")
}

func getCurrentWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
