package lib

import (
	"testing"
	"io"
	"bytes"
)

func testBody(t *testing.T, expected string, body io.Reader) {
	b := bytes.NewBuffer(nil)
	io.Copy(b, body)
	if actual := b.String(); expected != actual {
		t.Error(expected, actual)
	}
}

func TestDummyFetcherCount(t *testing.T) {
	f := NewDummyFetcherCount([]string{
		"res1",
		"res2",
	})
	if res, err := f.Get(""); err != nil {
		t.Error(err)
	} else {
		testBody(t, "res1", res.Body)
	}
	if res, err := f.Get(""); err != nil {
		t.Error(err)
	} else {
		testBody(t, "res2", res.Body)
	}
}
