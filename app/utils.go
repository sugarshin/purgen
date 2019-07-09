package app

import (
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func makeRequest(url string, ch chan<- http.Response) {
	res, err := http.Get(url)
	if err != nil {
		ch <- http.Response{}
	} else {
		ch <- *res
	}
}

func isURI(pathname string) bool {
	u, err := url.Parse(pathname)
	if err != nil || u.Scheme == "" {
		return false
	}
	return true
}

func getResourcesFromReader(body io.Reader) ([]string, error) {
	document, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, err
	}

	resources := []string{}

	document.Find("img, script, link").Each(func(index int, element *goquery.Selection) {
		src, exists := element.Attr("src")
		if exists {
			resources = append(resources, src)
		}

		href, exists := element.Attr("href")
		if exists {
			resources = append(resources, href)
		}

		srcset, exists := element.Attr("srcset")
		if exists {
			srcs := strings.Split(srcset, ",")
			regex := regexp.MustCompile(`\s+[0-9a-zA-Z]+$`)
			for i := range srcs {
				srcs[i] = strings.TrimSpace(srcs[i])
				srcs[i] = regex.ReplaceAllString(srcs[i], "")
			}
			resources = append(resources, srcs...)
		}
	})

	return resources, nil
}
