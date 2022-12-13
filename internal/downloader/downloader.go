package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"privateterraformregistry/internal/modules"
)

type Downloader interface {
	Download(w http.ResponseWriter, r *http.Request, m modules.Module) error
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

func (downloader *downloader) Download(w http.ResponseWriter, r *http.Request, m modules.Module) error {
	w.Header().Set("Content-Disposition", "attachment; filename=FILE_X")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))

	var fileName = fmt.Sprintf("%s.%s.%s.%s.tar.gz", m.Namespace, m.Name, m.System, m.Version)
	var filePath = downloader.config.DataDir + "/" + fileName
	fileToCopy, err := os.OpenFile(filePath, os.O_RDONLY, 0666)

	if err != nil {
		return err
	}

	defer fileToCopy.Close()
	_, err = io.Copy(w, fileToCopy)

	return err
}
