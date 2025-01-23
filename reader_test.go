package clireader_test

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"
	"testing"

	clireader "github.com/adamvduke/cli-reader"
	"github.com/adamvduke/go-collect"
)

type inOutPair struct {
	in  string
	out string
}

type testCase struct {
	pairs   []inOutPair
	bufSize int
}

var pairs = []inOutPair{
	{
		in:  "a",
		out: "a\n",
	},
	{
		in:  "a\n",
		out: "a\n",
	},
	{
		in:  "ab",
		out: "ab\n",
	},
	{
		in:  "ab\n",
		out: "ab\n",
	},
	{
		in:  "abc",
		out: "abc\n",
	},
	{
		in:  "abc\n",
		out: "abc\n",
	},
	{
		in:  "abcd",
		out: "abcd\n",
	},
	{
		in:  "abcd\n",
		out: "abcd\n",
	},
	{
		in:  "abcde",
		out: "abcde\n",
	},
	{
		in:  "abcde\n",
		out: "abcde\n",
	},
	{
		in:  "abcdef",
		out: "abcdef\n",
	},
	{
		in:  "abcdef\n",
		out: "abcdef\n",
	},
	{
		in:  "abcdefg",
		out: "abcdefg\n",
	},
	{
		in:  "abcdefg\n",
		out: "abcdefg\n",
	},
}

func TestReadAll(t *testing.T) {
	tc := testCase{pairs: pairs}
	reader := clireader.New(tc.inputs()...)
	for idx := range tc.inputs() {
		data, err := io.ReadAll(reader)
		if err != nil {
			// io.ReadAll does not consider io.EOF to be a reportable error. It's
			// used as a sentinel value to stop reading, so this error should be
			// nil.
			t.Errorf("unexpected error from io.ReadAll, error: %v", err)
		}
		got := string(data)
		want := tc.outputs()[idx]
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

const (
	maxBufSize = 20
)

func TestRead(t *testing.T) {
	cases := []testCase{}
	// generates test cases with varying buffer length and the set of pairs
	// defined above
	for i := 1; i <= maxBufSize; i++ {
		tc := testCase{pairs: pairs, bufSize: i}
		tc.bufSize = i
		cases = append(cases, tc)
	}

	for _, tc := range cases {
		t.Run(tc.Name(), func(t *testing.T) {
			reader := clireader.New(tc.inputs()...)
			for idx := range tc.inputs() {
				var err error
				var data []byte
				var n int
				// buf might be smaller than the current input, so keep
				// calling read until io.EOF
				for !errors.Is(err, io.EOF) {
					buf := make([]byte, tc.bufSize)
					n, err = reader.Read(buf)
					// only append up to n bytes to avoid indices in the buffer
					// that haven't been written
					data = append(data, buf[:n]...)
				}
				got := string(data)
				if got != tc.outputs()[idx] {
					t.Errorf("got: %s, want: %s", got, tc.outputs()[idx])
				}
			}
			// the reader should be exhausted at this point.
			count, err := reader.Read(make([]byte, tc.bufSize))
			if count != 0 {
				t.Errorf("unexpected number of bytes from Read got: %d, want: 0", count)
			}
			if !errors.Is(err, io.EOF) {
				t.Errorf("unexpected error from Read got: %v, want: %v", err, io.EOF)
			}
		})
	}
}

func (tc testCase) inputs() []string {
	return collect.Apply(tc.pairs, func(e inOutPair) string {
		return e.in
	})
}
func (tc testCase) outputs() []string {
	return collect.Apply(tc.pairs, func(e inOutPair) string {
		return e.out
	})
}

func (tc testCase) Name() string {
	inputs := []string{}
	for _, s := range tc.inputs() {
		next := strings.ReplaceAll(s, "\n", "")
		if !slices.Contains(inputs, next) {
			inputs = append(inputs, next)
		}
	}
	return fmt.Sprintf("%s_%d", strings.Join(inputs, "_"), tc.bufSize)
}
