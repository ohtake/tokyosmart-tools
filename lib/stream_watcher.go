package lib

import (
	"bufio"
	"net/http"
)

type StreamWatcher struct {
	endpoint  Endpoint
	lastFiles map[string]bool
	newFileCh chan string
}

func NewStreamWatcher(endpoint Endpoint, newFileCh chan string) *StreamWatcher {
	return &StreamWatcher{
		endpoint:  endpoint,
		newFileCh: newFileCh,
		lastFiles: make(map[string]bool),
	}
}

func (w *StreamWatcher) FetchList() error {
	resp, err := http.Get(w.endpoint.List())
	if err != nil {
		return err
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
