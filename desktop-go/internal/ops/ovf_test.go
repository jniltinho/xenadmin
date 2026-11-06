package ops

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
)

type failAfter struct {
	r   io.Reader
	err error
}

func (f *failAfter) Read(p []byte) (int, error) {
	n, err := f.r.Read(p)
	if errors.Is(err, io.EOF) {
		return n, f.err
	}
	return n, err
}

// TestExactReader: a body of its size goes; a longer one, a shorter one, or one whose decoder fails past its size fails.
func TestExactReader(t *testing.T) {
	read := func(r io.Reader, n int64) (string, error) {
		b, err := io.ReadAll(&exactReader{r: r, left: n})
		return string(b), err
	}
	if got, err := read(strings.NewReader("abcdef"), 6); err != nil || got != "abcdef" {
		t.Fatalf("%q %v", got, err)
	}
	for name, r := range map[string]io.Reader{"longer": strings.NewReader("abcdefg"), "shorter": strings.NewReader("abc"),
		"decoder":   &failAfter{strings.NewReader("abcdef"), errors.New("a grain out of order")},
		"truncated": &failAfter{strings.NewReader("abcdef"), fmt.Errorf("the disk is truncated: %w", io.EOF)}} {
		if _, err := read(r, 6); err == nil {
			t.Errorf("%s: no error", name)
		}
	}
}
