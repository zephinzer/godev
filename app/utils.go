package main

import (
	"os"
	"strings"
)

type commaSeparatedStringArray []string

func (sa *commaSeparatedStringArray) Set(item string) error {
	*sa = append(*sa, strings.Split(item, ",")...)
	return nil
}

func (sa *commaSeparatedStringArray) String() string {
	return strings.Join(*sa, ",")
}

func getCurrentWorkingDirectory() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}
