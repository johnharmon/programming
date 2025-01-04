package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type ProgressBar []byte

func (p ProgressBar) Render(progress int) {
	progressByte := byte(0174) //Pipe
	emptyByte := byte(040)     // Space
	for i := 0; i < 100; i++ {
		if i < progress {
			p[i] = progressByte
		} else if i >= progress && i < 99 {
			p[i] = emptyByte
		} else {
			p[i] = progressByte
		}
	}
}

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
		fmt.Println("Downloading with progress indicator")
		downloadFile, err := os.OpenFile("./download", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		downloadReader := urlDownloadResponse.Body
		if err != nil {
			return fmt.Errorf("error opening speficied file to write download to: %w", err)
		}
		defer downloadFile.Close()
		buf := make([]byte, chunkSize)
		p := make(ProgressBar, 100)
		for {
			bytesRead, err := downloadReader.Read(buf)
			if bytesRead > 0 {
				contentRead = contentRead + int64(bytesRead)
				percentDone = int64(float64(contentRead) / float64(contentLength) * 100)
				p.Render(int(percentDone))
				downloadFile.Write(buf[:bytesRead])
				if err == nil {
					fmt.Printf("%d%% %s\r", percentDone, string(p))
				}
				if err != nil && (err == io.EOF) {
					fmt.Printf("%d%% %s", percentDone, string(p))
					return nil
				} else if err != nil {
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
	//progressBar := make(ProgressBar, 100)

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
		} else {
			fmt.Println("")

		}
	}
}
