package xfile

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GetCurrentDirectory ...
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

// IsFileExists ...
func IsFileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func DownloadPictureBytes(url string, timeout time.Duration) (res []byte, err error) {
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer func() { _ = resp.Body.Close() }()
	fileByte := make([]byte, 0)
	fileByte, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return fileByte, nil

}
