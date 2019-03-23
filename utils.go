package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// ConfigCommaDelimitedString holds an array of strings to enable
// a single flag to take in multiple values via comma delimitation
type ConfigCommaDelimitedString []string

// Set adds a comma-delimited item to the main object
func (ccds *ConfigCommaDelimitedString) Set(item string) error {
	*ccds = append(*ccds, strings.Split(item, ",")...)
	return nil
}

// String representation
func (ccds *ConfigCommaDelimitedString) String() string {
	return strings.Join(*ccds, ",")
}

// ConfigMultiflagString holds an array of strings from a single
// flag that's been specified multiple times
type ConfigMultiflagString []string

// Set adds an item to the main object
func (cmfs *ConfigMultiflagString) Set(item string) error {
	*cmfs = append(*cmfs, item)
	return nil
}

// String representation
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

const confirmationTrueCanonical = "y"

var confirmationTrue = []string{confirmationTrueCanonical, "yes", "yupp", "yeah", "yea", "ok", "okay"}

const confirmationFalseCanonical = "n"

var confirmationFalse = []string{confirmationFalseCanonical, "no", "nope", "nah", "neh", "stop", "dont"}

func confirm(reader *bufio.Reader, question string, byDefault bool, retryText ...string) bool {
	var options string
	if byDefault {
		options = fmt.Sprintf("%s/%s", strings.ToUpper(confirmationTrueCanonical), confirmationFalseCanonical)
	} else {
		options = fmt.Sprintf("%s/%s", confirmationTrueCanonical, strings.ToUpper(confirmationFalseCanonical))
	}
	fmt.Printf("%s [%s]: ", question, options)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if len(userInput) < 2 {
		return byDefault
	}
	content := strings.Trim(
		strings.ToLower(userInput),
		" \r\n.,;",
	)
	confirmation := false
	if sliceContainsString(confirmationTrue, content) {
		confirmation = true
	} else if sliceContainsString(confirmationFalse, content) {
		confirmation = false
	} else if len(retryText) > 0 {
		fmt.Println(retryText[0])
		confirmation = confirm(reader, question, byDefault, retryText...)
	}
	return confirmation
}

func directoryExists(pathToDirectory string) bool {
	fileInfo, err := os.Lstat(pathToDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	if fileInfo.IsDir() {
		return true
	}
	return false
}

func fileExists(pathToFile string) bool {
	fileInfo, err := os.Lstat(pathToFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		panic(err)
	}
	if fileInfo.IsDir() {
		return false
	}
	return true
}

func sliceContainsString(slice []string, search string) bool {
	for _, sliceItem := range slice {
		if search == sliceItem {
			return true
		}
	}
	return false
}
