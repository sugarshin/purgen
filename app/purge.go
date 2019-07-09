package app

import (
	"net/http"
	"net/url"
	"path"
	"regexp"
	"time"
)

// Result : Purge result
type Result struct {
	URL    string
	Status int
}

func purge(resource string, baseURL string) Result {
	resourceURL := ""
	if path.IsAbs(resource) {
		base, err := url.Parse(baseURL)
		if err != nil {
			return Result{URL: "", Status: 0}
		}
		u, err := url.Parse(resource)
		if err != nil {
			return Result{URL: "", Status: 0}
		}
		resourceURL = base.ResolveReference(u).String()
	} else if isURI(resource) {
		resourceURL = resource
	} else {
		ext := path.Ext(baseURL)
		if ext != "" {
			dir, _ := path.Split(baseURL)
			resourceURL = regexp.MustCompile(`/$`).ReplaceAllString(dir, "") + "/" + path.Clean(resource)
		} else {
			resourceURL = regexp.MustCompile(`/$`).ReplaceAllString(baseURL, "") + "/" + path.Clean(resource)
		}
	}

	req, err := http.NewRequest("PURGE", resourceURL, nil)

	if err != nil {
		return Result{URL: resourceURL, Status: 0}
	}

	client := &http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	ch := make(chan http.Response)
	go makeClientRequest(client, req, ch)
	res := <-ch

	return Result{URL: resourceURL, Status: res.StatusCode}
}
