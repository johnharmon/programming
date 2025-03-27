package main

import "net/http"

type Response struct {
	StatusCode int           `json:"status_code"`
	Body       string        `json:"body"`
	Errors     []string      `json:"errors"`
	Request    *http.Request `json:"request"`
}
