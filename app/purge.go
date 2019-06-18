package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
)

// Result : Purge result
type Result struct {
	URL string
	Status int
}

func purge(imgSrc string, baseURL string) Result {
	imgSrcURL := ""
	_, err := url.ParseRequestURI(imgSrc)
	if err != nil { // is absolute path
		imgSrcURL = imgSrc
	} else {
		parsedBaseURL, _ := url.ParseRequestURI(baseURL)
		if err != nil {
			return Result{URL: "", Status:0}
		}
		parsedBaseURL.Path = path.Join(parsedBaseURL.Path, imgSrc)
		imgSrcURL = parsedBaseURL.String()
	}

	client := &http.Client{}
	req, err := http.NewRequest("PURGE", imgSrcURL, nil)
	if err != nil {
		log.Fatal(err)
		return Result{URL: imgSrcURL, Status:0}
	}
	res, err := client.Do(req)
	fmt.Println(res.StatusCode)
	if err != nil {
		log.Fatal(res, err)
		return Result{URL: imgSrcURL, Status:0}
	}

	return Result{URL: imgSrcURL, Status: res.StatusCode}
}
