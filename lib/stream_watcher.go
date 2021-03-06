package lib

import (
	"bufio"
	"fmt"
)

type StreamWatcher struct {
	endpoint  Endpoint
	lastFiles map[string]bool
	newFileCh chan string
	fetcher   Fetcher
}

func NewStreamWatcher(endpoint Endpoint, newFileCh chan string) *StreamWatcher {
	return NewStreamWatcherF(endpoint, newFileCh, NewDefaultFetcher())
}
func NewStreamWatcherF(endpoint Endpoint, newFileCh chan string, fetcher Fetcher) *StreamWatcher {
	return &StreamWatcher{
		endpoint:  endpoint,
		newFileCh: newFileCh,
		lastFiles: make(map[string]bool),
		fetcher:   fetcher,
	}
}

func (w *StreamWatcher) FetchList() error {
	uri := w.endpoint.List()
	resp, err := w.fetcher.Get(uri)
	if err != nil {
		return err
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response %d from %s", resp.StatusCode, uri)
	}
	files := make(map[string]bool)
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	var newFiles []string
	for scanner.Scan() {
		text := scanner.Text()
		files[text] = true
		seen := w.lastFiles[text]
		if !seen {
			newFiles = append(newFiles, text)
		}
	}
	for i := len(newFiles) - 1; i >= 0; i-- { // oldest file should be downloaded first
		w.newFileCh <- newFiles[i]
	}
	w.lastFiles = files
	return nil
}
