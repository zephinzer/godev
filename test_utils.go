package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
)

func createFile(t *testing.T, pathToFile string) {
	testFileCreation, err := os.Create(pathToFile)
	assert.Nilf(
		t,
		err,
		"unexpected error when creating file at '%s': %s",
		pathToFile,
		err,
	)
	defer testFileCreation.Close()
}

type MockCommand struct {
	Command
}

func ensureBoolFlag(t *testing.T, flag cli.Flag, matches string) {
	assert.IsType(t, cli.BoolFlag{}, flag)
	boolFlag := cli.BoolFlag(flag.(cli.BoolFlag))
	assert.Regexp(t, regexp.MustCompile(matches), boolFlag.Name)
}

func ensureDurationFlag(t *testing.T, flag cli.Flag, matches string) {
	assert.IsType(t, cli.DurationFlag{}, flag)
	durationFlag := cli.DurationFlag(flag.(cli.DurationFlag))
	assert.Regexp(t, regexp.MustCompile(matches), durationFlag.Name)
}

func ensureStringFlag(t *testing.T, flag cli.Flag, matches string) {
	assert.IsType(t, cli.StringFlag{}, flag)
	stringFlag := cli.StringFlag(flag.(cli.StringFlag))
	assert.Regexp(t, regexp.MustCompile(matches), stringFlag.Name)
}

func ensureStringSliceFlag(t *testing.T, flag cli.Flag, matches string) {
	assert.IsType(t, cli.StringSliceFlag{}, flag)
	stringSliceFlag := cli.StringSliceFlag(flag.(cli.StringSliceFlag))
	assert.Regexp(t, regexp.MustCompile(matches), stringSliceFlag.Name)
}

func ensureCLICommand(t *testing.T, command cli.Command, expectedName string, expectedAlias string, expectedFlags []cli.Flag) {
	assert.NotNil(t, command.Action)
	assert.Contains(t, command.Aliases, expectedAlias)
	assert.NotNil(t, command.Description)
	assert.Equal(t, expectedFlags, command.Flags)
	assert.Equal(t, expectedName, command.Name)
	assert.NotNil(t, command.Usage)
}

func ensureCLIFlags(t *testing.T, expectedFlags []string, actualFlags []cli.Flag) {
	matchedFlagsCount := 0
	for _, flag := range actualFlags {
		for _, expected := range expectedFlags {
			if strings.Contains(flag.GetName(), expected) {
				matchedFlagsCount++
				break
			}
		}
	}
	assert.Equal(t, matchedFlagsCount, len(expectedFlags))

}

func mockCommand(application string, arguments []string, logOutput *bytes.Buffer) *Command {
	command := &Command{
		id: fmt.Sprintf("%s%v", application, arguments),
		config: &CommandConfig{
			Application: application,
			Arguments:   arguments,
		},
		logger: InitLogger(&LoggerConfig{
			Name:   application,
			Format: "production",
			Level:  "trace",
		}),
		signal: make(chan os.Signal),
	}
	command.logger.SetOutput(logOutput)
	return command
}

func expectError(t *testing.T) func() {
	return func() {
		assert.NotNil(t, recover(), "expected an error but none was panicked")
	}
}

func expectNoError(t *testing.T) func() {
	return func() {
		err := recover()
		assert.Nilf(t, err, "expected no errors but '%s' was panicked", err)
	}
}

func removeFile(t *testing.T, pathToFile string) {
	err := os.Remove(pathToFile)
	assert.Nilf(
		t,
		err,
		"unexpected error when removing file at '%s': %s",
		pathToFile,
		err,
	)
}

func removeDir(t *testing.T, pathToDirectory string) error {
	directoryListing, err := ioutil.ReadDir(pathToDirectory)
	assert.Nil(t, err)
	for i := 0; i < len(directoryListing); i++ {
		listing := directoryListing[i]
		fullPath := path.Join(pathToDirectory, listing.Name())
		listingInfo, err := os.Lstat(fullPath)
		assert.Nil(t, err)
		if listingInfo.IsDir() {
			t.Logf("removing dir %s\n", fullPath)
			err = removeDir(t, fullPath)
		} else {
			err = os.Remove(fullPath)
			assert.Nil(t, err)
			t.Logf("removed file %s\n", fullPath)
		}
	}
	err = os.RemoveAll(pathToDirectory)
	return err
}
