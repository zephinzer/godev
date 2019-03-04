package main

import (
	"fmt"
	"os"
	"path"
)

type Initialiser interface {
	Check() bool
	Confirm() bool
	Handle(...bool) error
}

const InitialiserRetryText = "godev> sorry, i didn't get that"
const InitialiserFileSkipFormat = "godev> skipping '%s' - already exists"
const InitialiserDirSkipFormat = "godev> skipping - %s"

type FileInitialiser struct {
	Path     string
	Data     []byte
	Question string
}

func (fi FileInitialiser) Check() bool {
	return fileExists(fi.Path)
}

func (fi FileInitialiser) Confirm() bool {
	return confirm(
		Color("white", "godev> "+fi.Question),
		false,
		Color("bold", Color("red", InitialiserRetryText)),
	)
}

func (fi FileInitialiser) Handle(skip ...bool) error {
	if len(skip) > 0 && skip[0] {
		fmt.Println(
			Color("gray",
				fmt.Sprintf(InitialiserFileSkipFormat, path.Base(fi.Path)),
			),
		)
		return nil
	}
	filePath := path.Join(fi.Path)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, err = file.Write(fi.Data)
	if err != nil {
		return err
	}
	return nil
}

type DirInitialiser struct {
	Path        string
	Initialiser func() error
	Question    string
	Skip        string
}

func (di DirInitialiser) Check() bool {
	return directoryExists(di.Path)
}

func (di DirInitialiser) Confirm() bool {
	return confirm(
		Color("white", "godev> "+di.Question),
		false,
		InitialiserRetryText,
	)
}

func (di DirInitialiser) Handle(skip ...bool) error {
	if len(skip) > 0 && skip[0] {
		fmt.Println(
			Color("gray",
				fmt.Sprintf(InitialiserDirSkipFormat, di.Skip),
			),
		)
		return nil
	}
	return di.Initialiser()
}
