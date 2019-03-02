package main

import (
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)

var waitGroup sync.WaitGroup

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(path.Join(cwd, "../example/bin/app"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	wait := make(chan error)
	go func() {
		wait <- cmd.Run()
	}()
	select {
	case err := <-wait:
		if err != nil {
			log.Println(err)
		}
		log.Println(cmd.ProcessState)
	}
}
