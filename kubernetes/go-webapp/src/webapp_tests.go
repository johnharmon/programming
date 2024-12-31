package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type TokenTest struct {
	Name         string
	Url          string
	Method       string
	ExpectErrors bool
	Body         io.Reader
}

func MakeRequest(tt *TokenTest) (req *http.Request, reqErr error) {
	body := tt.Body
	req, err := http.NewRequest(tt.Method, tt.Url, body)
	if err != nil {
		reqErr = fmt.Errorf("error creating requesst: %+v", err)
	}
	return req, reqErr
}

var tokenTests = []TokenTest{
	{
		"Get-Token",
		"http://localhost:8080/jwt/token/get",
		"GET",
		false,
		nil,
	},
	{
		"Validate-token",
		"http://localhost:8080/jwt/token/validate",
		"GET",
		false,
		&bytes.Buffer{},
	},
}

func RequestTokenForClient(url string) (client *http.Client, reqErr error) {
	client = &http.Client{}
	_, err := client.Get(url)
	if err != nil {
		reqErr = fmt.Errorf("Error requesting token: %+v", err)
	}
	return client, reqErr
}
func TestTokens(t *testing.T) {
	client, err := RequestTokenForClient("http://localhost:8080/jwt/token/get")
	if err != nil {
		t.Errorf("error getting our token: %+v", err)
	}
	for _, test := range tokenTests {
		req, err := MakeRequest(&test)

	}
	return
}
