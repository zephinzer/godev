package main

import (
	"bufio"
)

// Initialiser is the interface for all file/directory bootstrapping
// operations
type Initialiser interface {
	Check() bool
	Confirm(*bufio.Reader) bool
	GetKey() string
	Handle(...bool) error
}

const initialiserRetryText = "godev> sorry, i didn't get that"
