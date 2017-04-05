package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

type TSDownloader struct {
	outputDir string
	endpoint  Endpoint
	newFileCh chan string
}

func NewTSDownloader(outputDir string, endpoint Endpoint, newFileCh chan string) TSDownloader {
	return TSDownloader{
		outputDir: outputDir,
		endpoint:  endpoint,
		newFileCh: newFileCh,
	}
}

func (d TSDownloader) Prepare() {
	os.Mkdir(d.outputDir, os.ModeDir)
}

func (d TSDownloader) Next() TSDownloaderResult {
	f := <-d.newFileCh
	localFilePath := path.Join(d.outputDir, TrimSerial(f))
	_, err := os.Stat(localFilePath)
	if err == nil {
		return newTSDownloaderResultSuccess(f, "skipped")
	}
	uri := d.endpoint.TS(f)
	resp, err := http.Get(uri)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	defer resp.Body.Close()
	file, err := os.Create(localFilePath)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	return newTSDownloaderResultSuccess(f, "downloaded")
}

type TSDownloaderResult struct {
	TsFile  string
	Success bool
	Action  string
	Error   error
}

func newTSDownloaderResultSuccess(ts string, action string) TSDownloaderResult {
	return TSDownloaderResult{
		TsFile:  ts,
		Success: true,
		Action:  action,
	}
}

func newTSDownloaderResultError(ts string, err error) TSDownloaderResult {
	return TSDownloaderResult{
		TsFile:  ts,
		Success: false,
		Error:   err,
	}
}
