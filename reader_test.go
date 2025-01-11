package clireader_test

import (
	"errors"
	"io"
	"testing"

	clireader "github.com/adamvduke/cli-reader"
)

func TestCLIReader_ReadAll(t *testing.T) {
	inputs := []string{"one", "two", "three\n"}
	wants := []string{"one\n", "two\n", "three\n"}
	reader := clireader.New(inputs...)
	for idx := range inputs {
		data, err := io.ReadAll(reader)
		if err != nil {
			// io.ReadAll does not consider io.EOF to be a reportable error. It's
			// used as a sentinel value to stop reading, so this error should be
			// nil.
			t.Errorf("unexpected error from io.ReadAll, error: %v", err)
		}
		got := string(data)
		want := wants[idx]
		if got != want {
			t.Errorf("got: %s, want: %s", got, want)
		}
	}
	// the reader should be exhausted at this point.
	data, err := io.ReadAll(reader)
	if err != nil {
		// io.ReadAll does not consider io.EOF to be a reportable error. It's
		// used as a sentinel value to stop reading, so this error should be
		// nil.
		t.Errorf("unexpected error from io.ReadAll, error: %v", err)
	}
	if len(data) != 0 {
		t.Errorf("unexpected number of bytes from io.ReadAll length: %d", data)
	}
}

func TestCLIReader_Read(t *testing.T) {
	inputs := []string{"one", "two", "three\n"}
	wants := []string{"one\n", "two\n", "three\n"}
	reader := clireader.New(inputs...)
	for idx := range inputs {
		var err error
		var data []byte
		// intentionally use a tiny buffer and read until io.EOF while
		// appending each result of calling Read
		for !errors.Is(err, io.EOF) {
			buf := make([]byte, 2)
			_, err = reader.Read(buf)
			data = append(data, buf...)
		}
		got := string(data)
		if got != wants[idx] {
			t.Errorf("got: %s, want: %s", got, wants[idx])
		}
	}
	// the reader should be exhausted at this point.
	count, err := reader.Read(make([]byte, 2))
	if count != 0 {
		t.Errorf("unexpected number of bytes from Read got: %d, want: 0", count)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("unexpected error from Read got: %v, want: %v", err, io.EOF)
	}
}
