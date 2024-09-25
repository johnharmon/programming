package main

import (
	"fmt"
	"io"
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
		return contentLength, nil
	}
}

func getUrl(url string, contentLength int64) (err error) {
	var useProgressBar bool
	chunkSize := 4096
	var contentRead int64
	var percentDone int64
	fmt.Printf("Content-Length: %d\n", contentLength)
	if contentLength == -1 {
		useProgressBar = false
	} else {
		useProgressBar = true
	}
	urlDownloadResponse, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error downloading the specified resource: %w", err)
	}
	if useProgressBar {
		//progressBar := make([]byte, 100)
		fmt.Println("Downloading with progress indicator")
		//progressBarTicSize = contentLength / 10000
		// "\0124"
		downloadFile, err := os.Open("./download")
		downloadReader := urlDownloadResponse.Body
		if err != nil {
			return fmt.Errorf("error opening speficied file to write download to: %w", err)
		}
		defer downloadFile.Close()
		buf := make([]byte, chunkSize)
		for {
			bytesRead, err := downloadReader.Read(buf)
			//fmt.Printf("Read chunk\n")
			if bytesRead > 0 {
				contentRead = contentRead + int64(bytesRead)
				percentDone = int64(float64(contentRead) / float64(contentLength) * 100)
				downloadFile.Write(buf)
				fmt.Printf("%d%% Done\r", percentDone)
				if err != nil && (err == io.EOF) {
					return nil
				} else if err != nil {
					//fmt.Printf("Unknown error, %+v", err)
					return fmt.Errorf("error: %w", err)
				}
			}
			if err != nil && err == io.EOF {
				return nil
			} else if err != nil && err != io.EOF {
				return fmt.Errorf("error reading from stream: %w", err)
			} else if err != nil {
				return fmt.Errorf("unknown error: %w", err)
			}
		}
	}
	return nil
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
		err := getUrl(url, contentLength)
		if err != nil {
			//log.Fatal("Error: %s\n", err)
			log.Fatalf("Error: %+v\n", err)
		}
	}
}
