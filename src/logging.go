package goule

import (
	"io"
	"io/ioutil"
)

// createLogStdout takes an already-locked executable and creates a stream for
// its standard output.
func createLogStdout(exc *Executable) (io.Writer, error) {
	// TODO: here, open a file etc.
	return ioutil.Discard, nil
}

// createLogStderr takes an already-locked executable and creates a steram for
// its standard error.
func createLogStderr(exc *Executable) (io.Writer, error) {
	// TODO: here, open a file etc.
	return ioutil.Discard, nil
}