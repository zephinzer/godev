package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Join(cwd, "/.output"), []byte("It works!"), 0644); err != nil {
		panic(err)
	}
	logrus.Info("It works!")
}
