package lib

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

type errorReader struct {
}

func (errorReader) Read(p []byte) (int, error) {
	return 0, errors.New("always error")
}

type dummyTSWriter struct {
	files      map[string]bool
	customOpen func(string) (handled bool, _ io.WriteCloser, _ error)
}

func newDummyTSWriter(customOpen func(string) (bool, io.WriteCloser, error)) dummyTSWriter {
	return dummyTSWriter{
		files:      make(map[string]bool),
		customOpen: customOpen,
	}
}
func (w dummyTSWriter) Prepare() {
}
func (w dummyTSWriter) HasTS(tsFile string) bool {
	return w.files[tsFile]
}
func (w dummyTSWriter) Open(tsFile string) (io.WriteCloser, error) {
	if nil != w.customOpen {
		handled, res1, res2 := w.customOpen(tsFile)
		if handled {
			if res2 != nil {
				w.files[tsFile] = true
			}
			return res1, res2
		}
	}
	w.files[tsFile] = true
	return nopWriteCloser{}, nil
}
func (w dummyTSWriter) Remove(tsFile string) error {
	delete(w.files, tsFile)
	return nil
}
func (w dummyTSWriter) numberOfFiles() int {
	return len(w.files)
}

type fetcherRaisesErrorWhileDownloading struct {
	body string
}

func (f *fetcherRaisesErrorWhileDownloading) Get(url string) (*http.Response, error) {
	reader1 := strings.NewReader(f.body)
	reader2 := errorReader{}
	reader := ioutil.NopCloser(io.MultiReader(reader1, reader2))
	return &http.Response{
		StatusCode: 200,
		Body:       reader,
	}, nil
}

func TestTSDownloader_ShoudNotDownloadOrOpenWhenFileExists(t *testing.T) {
	w := newDummyTSWriter(nil)
	w.Open("exists.ts")
	if w.numberOfFiles() != 1 {
		t.Error(w.numberOfFiles())
	}
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{"ts"}), w)
	ch <- "exists.ts"
	if result := d.Next(); !result.Success || result.Action != "skipped" {
		t.Error(result)
	}
	if w.numberOfFiles() != 1 {
		t.Error(w.numberOfFiles())
	}
	ch <- "exists-not.ts"
	if result := d.Next(); !result.Success || result.Action != "downloaded" {
		t.Error(result)
	}
	if w.numberOfFiles() != 2 {
		t.Error(w.numberOfFiles())
	}
}

func TestTSDownloader_ShoudReturnFailureIfErrorOnSave(t *testing.T) {
	w := newDummyTSWriter(func(string) (bool, io.WriteCloser, error) {
		return true, errorWriteCloser{}, nil
	})
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{"non-empty to be written to file"}), w)
	ch <- "file.ts"
	if result := d.Next(); result.Success {
		t.Error(result)
	}
	if w.numberOfFiles() != 0 {
		t.Error(w.numberOfFiles())
	}
}

func TestTSDownloader_ShoudReturnFailureIfFetcherFailsAtFirst(t *testing.T) {
	w := newDummyTSWriter(nil)
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, NewDummyFetcherCount([]string{}), w)
	ch <- "file.ts"
	if result := d.Next(); result.Success {
		t.Error(result)
	}
	if w.numberOfFiles() != 0 {
		t.Error(w.numberOfFiles())
	}
}

func TestTSDownloader_ShoudReturnFailureIfFetcherFailsWhileDonwloading(t *testing.T) {
	w := newDummyTSWriter(nil)
	ch := make(chan string, 10)
	d := NewTSDownloaderFW(NewEndpoint("v-low-tokyo1", "2865"), ch, &fetcherRaisesErrorWhileDownloading{"some data before error"}, w)
	ch <- "file.ts"
	if result := d.Next(); result.Success {
		t.Error(result)
	}
	if w.numberOfFiles() != 0 {
		t.Error(w.numberOfFiles())
	}
}
