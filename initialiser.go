package main

import (
	"bufio"
)

type Initialiser interface {
	Check() bool
	Confirm(*bufio.Reader) bool
	GetKey() string
	Handle(...bool) error
}

const InitialiserRetryText = "godev> sorry, i didn't get that"
