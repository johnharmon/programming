package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"testing"
)

var tokenTests = []TokenTest{
	{
		"Get-Token",
		"http://localhost:8080/jwt/token/get",
		"GET",
		false,
		nil,
		nil,
		nil,
	},
	{
		"Validate-token",
		"http://localhost:8080/jwt/token/validate",
		"GET",
		false,
		&bytes.Buffer{},
		nil,
		nil,
	},
}

type TokenTest struct {
	Name         string
	Url          string
	Method       string
	ExpectErrors bool
	Body         io.Reader
	Request      *http.Request
	Client       *http.Client
}

type TokenTestConfig struct {
	Name         string
	Url          string
	Method       string
	ExpectErrors bool
	Body         io.Reader
	Request      *http.Request
	Client       *http.Client
}

func MakeRequest(tt *TokenTest) (req *http.Request, reqErr error) {
	body := tt.Body
	req, err := http.NewRequest(tt.Method, tt.Url, body)
	if err != nil {
		reqErr = fmt.Errorf("error creating requesst: %+v", err)
	}
	return req, reqErr
}

func (tt *TokenTest) MakeRequest() (reqErr error) {
	body := tt.Body
	req, err := http.NewRequest(tt.Method, tt.Url, body)
	if err != nil {
		reqErr = fmt.Errorf("error creating requesst: %+v", err)
	}
	tt.Request = req
	return reqErr
}

func RequestTokenForClient(url string) (client *http.Client, reqErr error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Printf("Error creating cookiejar: %+v", err)
		reqErr = errors.Join(reqErr, fmt.Errorf("error creating cookiejar: %+v", err))
	}
	client = &http.Client{
		Jar: jar,
	}
	_, err = client.Get(url)
	//fmt.Printf("Client: %+v\n", client)
	if err != nil {
		reqErr = errors.Join(reqErr, fmt.Errorf("error requesting token: %+v", err))
	}
	return client, reqErr
}

func testToken(test TokenTest, t *testing.T) {
	req, reqErr := MakeRequest(&test)
	if reqErr != nil {
		t.Fatalf("error generating request: %+v", reqErr)
	}
	resp, doErr := test.Client.Do(req)
	if doErr != nil {
		t.Errorf("error making request: %+v", doErr)
	}
	if resp != nil {
		body, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			t.Errorf("error reading repsonse body: %+v", bodyErr)
		}
		t.Logf("response body: %+v", string(body))
		//fmt.Printf("t: %v\n", t)
	}
}

func TestTokens(t *testing.T) {
	client, err := RequestTokenForClient("http://localhost:8080/jwt/token/get")
	if err != nil {
		t.Fatalf("error getting our token: %+v", err)
	}
	for _, test := range tokenTests {
		test.Client = client
		//fmt.Printf("Testing case: %+v\n", test)
		reqUrl, err := url.Parse(test.Url)
		fmt.Printf("URL: %+v\n", reqUrl)
		testToken(test, t)
		if err != nil {
			t.Errorf("Error generating url from: %s\nError was: %+v", test.Url, err)
		}
		//fmt.Printf("Client info: %+v\n", test.Client)
		//cookies := client.Jar.Cookies(reqUrl)
		//fmt.Printf("%s\n", cookies)
		//fmt.Printf("Client Cookies: %+v", test.Client.Jar.Cookies(reqUrl))

	}
}
