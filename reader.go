// Package clireader provides an io.Reader intended to simulate newline delimited
// input from the command line.
package clireader

import (
	"io"
	"strings"
	"sync"
)

// Reader implements an io.Reader intended to simulate newline delimited
// input from the command line. A Reader is constructed by providing a
// list of strings. Each call to Read will return a string from the provided
// list, in order, with a newline appended, along with an error value of io.EOF
// to signal the end of the current input. Once all strings have been read,
// calling Read stops reading bytes and returns (0, io.EOF) to signal that the
// reader has been exhausted.
type Reader struct {
	inputs     [][]byte
	inputIndex int
	byteIndex  int
	mu         *sync.Mutex
}

// New returns a *Reader configured to return a string from the provided
// list, in order, with a newline appended, for each call to Read.
func New(inputs ...string) *Reader {
	reader := &Reader{inputs: make([][]byte, len(inputs)), mu: &sync.Mutex{}}
	for idx, input := range inputs {
		// make sure all the inputs end with a newline so it can be used as a
		// delimiter when reading
		if !strings.HasSuffix(input, "\n") {
			input += "\n"
		}
		reader.inputs[idx] = []byte(input)
	}

	return reader
}

// Read copies up to len(buf) bytes from the Reader into buf, and returns the number
// of bytes copied along with any error.
//
// If the number of bytes returned is greater than 0, and the returned error is nil,
// there are more bytes to be read from the current input.
//
// If the number of bytes returned is greater than 0, and the returned error is io.EOF
// the current input is complete, and the reader will be advanced to the next input.
//
// If the number of bytes returned is 0, and the returned error is io.EOF, the reader
// has been exhausted.
func (reader *Reader) Read(buf []byte) (int, error) {
	reader.mu.Lock()
	defer reader.mu.Unlock()

	// all inputs have been read, nothing to do
	if reader.inputIndex == len(reader.inputs) {
		return 0, io.EOF
	}

	// attempt to copy from reader.byteIndex to the end of the current input and
	// update reader.byteIndex based on the number of bytes read
	curr := reader.inputs[reader.inputIndex]
	readCount := copy(buf, curr[reader.byteIndex:])
	reader.byteIndex += readCount

	// if the current input is completely read, advance to the next and return io.EOF
	// to indicate the end of the input
	if reader.byteIndex == len(curr) {
		reader.inputIndex++
		reader.byteIndex = 0
		return readCount, io.EOF
	}

	// there's more data to read, so only return the number of bytes read and a nil error
	return readCount, nil
}
