package app

import (
	"net/http"
)

// MakeRequest : http request
func MakeRequest(url string, ch chan<- http.Response) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- http.Response{}
	} else {
		ch <- *resp
	}
}
