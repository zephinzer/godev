package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
)

// GitInitialiserConfig holds the configurations for the GitInitialiser
type GitInitialiserConfig struct {
	Path string
}

// InitGitInitialiser initialises the Git initialiser that assists in
// initialising a directory as a Git repository
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

// GitInitialiser assists in initialising a directory as a Git repository
type GitInitialiser struct {
	Key    string
	Path   string
	logger *Logger
}

// Check verifies that the path exists
func (gi *GitInitialiser) Check() bool {
	return directoryExists(path.Join(gi.Path, "/.git"))
}

// Confirm seeks advice from the user whether we should proceed with the
// Git repository initialisation
func (gi *GitInitialiser) Confirm(reader *bufio.Reader) bool {
	return confirm(
		reader,
		Color("white", fmt.Sprintf("godev> initialise git repository at '%s'?", gi.Path)),
		false,
		Color("bold", Color("red", initialiserRetryText)),
	)
}

// GetKey returns the key of this initialiser
func (gi *GitInitialiser) GetKey() string {
	return gi.Key
}

// Handle processes the initialiser (initialises the Git repository)
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
