package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func headUrl(url string) (contentLength int64, err error) {
	headResponse, err := http.Head(url)
	if err != nil {
		return -1, fmt.Errorf("error getting information about download: %w", err)
	} else {
		contentLength = headResponse.ContentLength
	}
}

func getUrl(url string, contentLength int64) (err error) {

}

func main() {
	var url string

	if len(os.Args) < 2 {
		url = "https://go.dev/dl/go1.23.1.darwin-amd64.pkg"
	} else {
		url = os.Args[1]
	}
	contentLength, err := headUrl(url)
	if err != nil {
		log.Fatal("Error fetching downloadable content: %w", err)
	} else if contentLength == -1 {
		fmt.Printf("Unknown content length, progress estimation will be... unreliable\n")
		fmt.Printf("Downloading %s\n", url)
		getUrl(url, contentLength)
	} else {
		fmt.Printf("Downloading %s\n", url)
	}

}
