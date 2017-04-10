package lib

import (
	"io"
)

type TSDownloader struct {
	endpoint  Endpoint
	newFileCh chan string
	fetcher   Fetcher
	writer    TSWriter
}

func NewTSDownloader(outputDir string, endpoint Endpoint, newFileCh chan string) TSDownloader {
	return NewTSDownloaderFW(endpoint, newFileCh, NewDefaultFetcher(), NewDefaultTSWriter(outputDir))
}

func NewTSDownloaderFW(endpoint Endpoint, newFileCh chan string, fetcher Fetcher, writer TSWriter) TSDownloader {
	return TSDownloader{
		endpoint:  endpoint,
		newFileCh: newFileCh,
		fetcher:   fetcher,
		writer:    writer,
	}
}

func (d TSDownloader) Prepare() {
	d.writer.Prepare()
}

func (d TSDownloader) Next() TSDownloaderResult {
	f := <-d.newFileCh
	if d.writer.HasTS(f) {
		return newTSDownloaderResultSuccess(f, "skipped")
	}
	uri := d.endpoint.TS(f)
	resp, err := d.fetcher.Get(uri)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	defer resp.Body.Close()
	out, err := d.writer.Open(f)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	incompleteDownload := true
	defer func() {
		out.Close()
		if incompleteDownload {
			d.writer.Remove(f)
		}
	}()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return newTSDownloaderResultError(f, err)
	}
	incompleteDownload = false
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
