package lib

import (
	"io"
	"os"
	"path"
)

type TSWriter interface {
	Prepare()
	HasTS(tsFile string) bool
	Open(tsFile string) (io.WriteCloser, error)
}

type defaultTSWriter struct {
	outputDir string
}

func NewDefaultTSWriter(outputDir string) TSWriter {
	return defaultTSWriter{
		outputDir: outputDir,
	}
}

func (w defaultTSWriter) Prepare() {
	os.Mkdir(w.outputDir, 0755)
}

func (w defaultTSWriter) HasTS(tsFile string) bool {
	localFilePath := path.Join(w.outputDir, TrimSerial(tsFile))
	_, err := os.Stat(localFilePath)
	return err == nil
}

func (w defaultTSWriter) Open(tsFile string) (io.WriteCloser, error) {
	localFilePath := path.Join(w.outputDir, TrimSerial(tsFile))
	return os.Create(localFilePath)
}
