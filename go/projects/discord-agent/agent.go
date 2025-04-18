package discordagent

import (
	"io"
	"net/http"
	u "net/url"
	"os"
)

func writeURLContent(resp *http.Response, f *os.File) {
	buf := make([]byte, 8192)
	for {
		b, err := resp.Body.Read(buf)
		if err != nil {
			if err != io.EOF {
			}
			break
		} else if b == 0 {
			break
		}
		f.Write(buf)
	}
}

func downloadUrl(url u.URL) (downloadErr error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}
	writeURLContent(resp)
	return downloadErr
}

func runDiscordBot(downloads chan u.URL) {
	return
}

func runDownloader(downloads chan u.URL) {
	for du := range downloads {
		go downloadUrl(du)
	}
	return
}

func main() {
	downloads := make(chan u.URL)
	go runDiscordBot(downloads)
	go runDownloader(downloads)

}
