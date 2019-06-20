package app

import (
	"fmt"
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

func purge(imgSrc string, baseURL string) Result {
	imgSrcURL := ""
	if path.IsAbs(imgSrc) {
		base, err := url.Parse(baseURL)
		if err != nil {
			return Result{URL: "", Status: 0}
		}
		u, err := url.Parse(imgSrc)
		if err != nil {
			return Result{URL: "", Status: 0}
		}
		imgSrcURL = base.ResolveReference(u).String()
	} else if IsURI(imgSrc) {
		imgSrcURL = imgSrc
	} else {
		ext := path.Ext(baseURL)
		if ext != "" {
			dir, _ := path.Split(baseURL)
			imgSrcURL = regexp.MustCompile(`/$`).ReplaceAllString(dir, "") + "/" + path.Clean(imgSrc)
		} else {
			imgSrcURL = regexp.MustCompile(`/$`).ReplaceAllString(baseURL, "") + "/" + path.Clean(imgSrc)
		}
	}

	client := &http.Client{}
	req, err := http.NewRequest("PURGE", imgSrcURL, nil)
	if err != nil {
		log.Fatal(err)
		return Result{URL: imgSrcURL, Status: 0}
	}
	res, err := client.Do(req)
	fmt.Println(res.StatusCode)
	if err != nil {
		log.Fatal(res, err)
		return Result{URL: imgSrcURL, Status: 0}
	}

	return Result{URL: imgSrcURL, Status: res.StatusCode}
}
