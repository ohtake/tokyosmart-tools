package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	newFileCh := make(chan string, 2*100)
	go func() {
		streamWather := NewStreamWatcher("https://smartcast.hs.llnwd.net/v-low-tokyo1/2865/2865.txt", newFileCh)
		for {
			err := streamWather.FetchList()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		downloader := NewDownloader("output", "https://smartcast.hs.llnwd.net/v-low-tokyo1/2865/", newFileCh)
		downloader.Prepare()
		for {
			err := downloader.Next()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}()
	select {}
}

type StreamWatcher struct {
	url       string
	lastFiles map[string]bool
	newFileCh chan string
}

func NewStreamWatcher(url string, newFileCh chan string) *StreamWatcher {
	return &StreamWatcher{
		url:       url,
		newFileCh: newFileCh,
		lastFiles: make(map[string]bool),
	}
}

func (w *StreamWatcher) FetchList() error {
	resp, err := http.Get(w.url)
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

type Downloader struct {
	outputDir string
	baseURI   string
	newFileCh chan string
}

func NewDownloader(outputDir string, baseURI string, newFileCh chan string) Downloader {
	return Downloader{
		outputDir: outputDir,
		baseURI:   baseURI,
		newFileCh: newFileCh,
	}
}

func (d Downloader) Prepare() {
	os.Mkdir(d.outputDir, os.ModeDir)
}

func (d Downloader) Next() error {
	f := <-d.newFileCh
	localFilePath := path.Join(d.outputDir, f[3:])
	_, err := os.Stat(localFilePath)
	if err == nil {
		fmt.Println("skipped " + f)
		return nil
	}
	uri := d.baseURI + f
	resp, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	file, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("downloaded " + f)
	return nil
}
