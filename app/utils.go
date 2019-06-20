package app

import (
	"net/http"
	"net/url"
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

// IsURI : is uri
func IsURI(pathname string) bool {
	u, err := url.Parse(pathname)
	if err != nil || u.Scheme == "" {
		return false
	}
	return true
}
