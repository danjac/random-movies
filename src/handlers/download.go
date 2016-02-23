package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type downloader interface {
	download(string, string, string) (string, error)
}

type imdbPosterDownloader struct{}

func (d *imdbPosterDownloader) download(staticDir, url, imdbID string) (string, error) {

	if url == "" {
		return "", nil
	}

	filename := fmt.Sprintf("%s.jpg", imdbID)

	imageDir := filepath.Join(staticDir, "images")
	if err := os.Mkdir(imageDir, 0777); err != nil && !os.IsExist(err) {
		return filename, err
	}

	imagePath := filepath.Join(imageDir, filename)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return filename, err
	}
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		return filename, err
	}

	// should probably check this
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filename, err
}
