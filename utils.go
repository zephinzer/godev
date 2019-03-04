package main

import (
	"bufio"
	"fmt"
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

type Confirmation *bool

const ConfirmationTrueCanonical = "y"

var ConfirmationTrue = []string{ConfirmationTrueCanonical, "yes", "yupp", "yeah", "yea", "ok", "okay"}

const ConfirmationFalseCanonical = "n"

var ConfirmationFalse = []string{ConfirmationFalseCanonical, "no", "nope", "nah", "neh", "stop", "dont"}

func confirm(question string, byDefault bool, retryText ...string) Confirmation {
	var options string
	if byDefault {
		options = fmt.Sprintf("%s/%s", strings.ToUpper(ConfirmationTrueCanonical), ConfirmationFalseCanonical)
	} else {
		options = fmt.Sprintf("%s/%s", ConfirmationTrueCanonical, strings.ToUpper(ConfirmationFalseCanonical))
	}
	fmt.Printf("%s [%s]: ", question, options)
	reader := bufio.NewReader(os.Stdin)
	userInput, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	if len(userInput) < 2 {
		return Confirmation(&byDefault)
	} else {
		content := strings.Trim(strings.ToLower(userInput), " \n")
		confirmation := true
		if sliceContainsString(ConfirmationTrue, content) {
			confirmation = true
		} else if sliceContainsString(ConfirmationFalse, content) {
			confirmation = false
		} else if len(retryText) > 0 {
			fmt.Println(retryText[0])
			confirmation = *confirm(question, byDefault, retryText...)
		} else {
			return nil
		}
		return &confirmation
	}
}

func directoryExists(pathToDirectory string) bool {
	fileInfo, err := os.Lstat(pathToDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			panic(err)
		}
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
		} else {
			panic(err)
		}
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
