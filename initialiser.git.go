package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type GitInitialiserConfig struct {
	Path string
}

func InitGitInitialiser(config *GitInitialiserConfig) *GitInitialiser {
	gi := &GitInitialiser{
		Key:  ".git",
		Path: config.Path,
		logger: InitLogger(&LoggerConfig{
			Format: "raw",
		}),
	}
	return gi
}

type GitInitialiser struct {
	Key    string
	Path   string
	logger *Logger
}

func (gi *GitInitialiser) Check() bool {
	return directoryExists(path.Join(gi.Path, "/.git"))
}

func (gi *GitInitialiser) Confirm(reader *bufio.Reader) bool {
	return confirm(
		reader,
		Color("white", fmt.Sprintf("godev> initialise git repository at '%s'?", gi.Path)),
		false,
		Color("bold", Color("red", InitialiserRetryText)),
	)
}

func (gi *GitInitialiser) GetKey() string {
	return gi.Key
}

func (gi *GitInitialiser) Handle(skip ...bool) error {
	if len(skip) > 0 && skip[0] {
		gi.logger.Info(
			Color("gray",
				fmt.Sprintf("godev> skipping git repository initialisation at '%s'", gi.Path),
			),
		)
		return nil
	}
	var err error
	if _, err = exec.LookPath("git"); err != nil {
		return err
	}
	cmd := exec.Command("git", "init")
	cmd.Dir = gi.Path
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	return err
}
