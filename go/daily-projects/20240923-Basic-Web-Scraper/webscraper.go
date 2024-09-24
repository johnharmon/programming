package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func queryUrl(url string) (urlBody []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error querying url: %w", err)
	} else {
		defer resp.Body.Close()
		urlBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading content: %w", err)
		}
		return urlBody, nil
	}
}

func analyzeUrlBody(urlBody []byte) ([][]byte, error) {
	linkRegex, err := regexp.Compile("<a .*?href=.*?(/>|</a>)")
	if err != nil {
		return nil, fmt.Errorf("error compiling regex: %w", err)
	}
	matches := linkRegex.FindAll(urlBody, -1)
	return matches, nil
}

func showMatches(matches [][]byte) {
	for index, match := range matches {
		fmt.Printf("%d: %s\n", index+1, string(match))
	}
}

func main() {
	var url string
	if len(os.Args) < 2 {
		url = "https://pkg.go.dev/std"
	} else {
		url = os.Args[1]
	}
	urlBody, err := queryUrl(url)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(1)
	}
	matches, err := analyzeUrlBody(urlBody)
	if err != nil {
		fmt.Printf("Error: %+v\n", err)
		os.Exit(2)
	}
	showMatches(matches)
}
