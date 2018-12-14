# Golang Development Image - Test
This directory contains a test application which is used to test the build/packaging process

Run `make test` from the root directory of this repository to invoke the building of this directory.

The application contained herein tests for:

- successful `go build`
- successful usage of external dependencies from the golang standard library
- successful usage of external dependencies from github
- successful running of a file-modifying function
- successful reloading of a file based change

# The Dockerfile
The Dockerfile here is a sample Dockerfile you can copy and paste into your own projects for use
