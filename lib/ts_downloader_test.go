package lib

import (
	"errors"
	"io"
	"testing"
)

type nopWriteCloser struct {
}

func (nopWriteCloser) Write(p []byte) (int, error) {
	return len(p), nil
}
func (nopWriteCloser) Close() error {
	return nil
}

type errorWriteCloser struct {
}

func (errorWriteCloser) Write(p []byte) (int, error) {
	return 0, errors.New("always error")
}
func (errorWriteCloser) Close() error {
	return nil
}

type dummyTSWriter struct {
	hasTS func(string) bool
	open  func(string) (io.WriteCloser, error)
}

func newDummyTSWriter(hasTS func(string) bool, open func(string) (io.WriteCloser, error)) TSWriter {
	return dummyTSWriter{
		hasTS: hasTS,
		open:  open,
	}
}
func (w dummyTSWriter) Prepare() {
}
func (w dummyTSWriter) HasTS(tsFile string) bool {
	return w.hasTS(tsFile)
}
func (w dummyTSWriter) Open(tsFile string) (io.WriteCloser, error) {
	return w.open(tsFile)
}

func TestTSDownloader_ShoudNotDownloadOrOpenWhenFileExists(t *testing.T) {
	openCallCount := 0
	w := newDummyTSWriter(
		func(f string) bool { return "exists.ts" == f },
		func(string) (io.WriteCloser, error) {
			openCallCount++
			return nopWriteCloser{}, nil
		},
	)
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{""}), w)
	ch <- "exists.ts"
	if result := d.Next(); !result.Success || result.Action != "skipped" {
		t.Error(result)
	}
	if openCallCount != 0 {
		t.Error(openCallCount)
	}
	ch <- "exists-not.ts"
	if result := d.Next(); !result.Success || result.Action != "downloaded" {
		t.Error(result)
	}
	if openCallCount != 1 {
		t.Error(openCallCount)
	}
}

func TestTSDownloader_ShoudReturnFailureIfErrorOnSave(t *testing.T) {
	w := newDummyTSWriter(
		func(f string) bool { return false },
		func(string) (io.WriteCloser, error) {
			return errorWriteCloser{}, nil
		},
	)
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{"non-empty to be written to file"}), w)
	ch <- "file.ts"
	if result := d.Next(); result.Success {
		t.Error(result)
	}
}

func TestTSDownloader_ShoudReturnFailureIfFetcherFails(t *testing.T) {
	openCallCount := 0
	w := newDummyTSWriter(
		func(f string) bool { return false },
		func(string) (io.WriteCloser, error) {
			openCallCount++
			return nopWriteCloser{}, nil
		},
	)
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{}), w)
	ch <- "file.ts"
	if result := d.Next(); result.Success {
		t.Error(result)
	}
	if openCallCount != 0 {
		t.Error(openCallCount)
	}
}
