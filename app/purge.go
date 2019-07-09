package app

import (
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
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

	client := http.Client{}
	req, err := http.NewRequest("PURGE", resourceURL, nil)
	if err != nil {
		log.Fatal(err)
		return Result{URL: resourceURL, Status: 0}
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return Result{URL: resourceURL, Status: 0}
	}

	return Result{URL: resourceURL, Status: res.StatusCode}
}
