package xauth

import (
	"fmt"
	"io"
	"net/http"
)

type HttpGetResponse struct {
	Body    []byte
	Headers http.Header
}

func HttpGet(client *http.Client, url string) (response HttpGetResponse, err error) {
	r, err := client.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	response = HttpGetResponse{body, r.Header}
	if r.StatusCode >= 300 {
		err = fmt.Errorf(string(response.Body))
		return
	}
	fmt.Printf("HTTP GET %s: %s %s \n", url, r.Status, string(response.Body))
	err = nil
	return
}
