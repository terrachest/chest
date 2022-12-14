package downloader

import (
	"io"
	"net/http"
	"os"
)

type Downloader interface {
	Download(w http.ResponseWriter, r *http.Request, fn string) error
}

func New(c *Config) Downloader {
	return &downloader{
		config: c,
	}
}

type Config struct {
	DataDir string
}

type downloader struct {
	config *Config
}

func (downloader *downloader) Download(w http.ResponseWriter, r *http.Request, fn string) error {
	w.Header().Set("Content-Disposition", "attachment; filename=FILE_X")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	var filePath = downloader.config.DataDir + "/" + fn
	fileToCopy, err := os.OpenFile(filePath, os.O_RDONLY, 0666)

	if err != nil {
		return err
	}

	defer fileToCopy.Close()
	_, err = io.Copy(w, fileToCopy)

	return err
}
