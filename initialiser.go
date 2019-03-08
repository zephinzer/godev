package main

import (
	"bufio"
)

type Initialiser interface {
	Check() bool
	Confirm(*bufio.Reader) bool
	Handle(...bool) error
}

const InitialiserRetryText = "godev> sorry, i didn't get that"
