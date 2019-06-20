package app

import (
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

func setupRoutes(engine *gin.Engine) *gin.Engine {
	engine.GET("/", handleIndex)
	engine.GET("/ping", handlePing)
	engine.POST("/purge", handlePurge)
	engine.NoRoute(hanldeNotFound)
	return engine
}

func handleIndex(c *gin.Context) {
	c.HTML(200, "index.tmpl", gin.H{})
}

func handlePurge(c *gin.Context) {
	formDataURL := c.PostForm("url")
	_, err := url.ParseRequestURI(formDataURL)
	if err != nil {
		c.HTML(400, "index.tmpl", gin.H{
			"error": "URL must be correct",
		})
		return
	}

	res := http.Response{}
	ch := make(chan http.Response)
	go MakeRequest(formDataURL, ch)
	res = <-ch
	if res.Body == nil {
		c.HTML(400, "index.tmpl", gin.H{
			"error": "Server IP address could not be found.",
		})
		return
	}

	defer res.Body.Close()

	document, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body", err)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"error": err,
		})
		return
	}

	results := []Result{}
	document.Find("img").Each(func(index int, element *goquery.Selection) {
		imgSources := []string{}
		imgSrc, exists := element.Attr("src")
		if exists {
			imgSources = append(imgSources, imgSrc)
		}

		srcset, exists := element.Attr("srcset")
		if exists {
			srcs := strings.Split(srcset, ",")
			regex := regexp.MustCompile(`\s+[0-9a-zA-Z]+$`)
			for i := range srcs {
				srcs[i] = strings.TrimSpace(srcs[i])
				srcs[i] = regex.ReplaceAllString(srcs[i], "")
			}
			imgSources = append(imgSources, srcs...)
		}

		for i := range imgSources {
			results = append(results, purge(imgSources[i], formDataURL))
		}
	})

	c.HTML(200, "index.tmpl", gin.H{
		"results": results,
	})
}

func handlePing(c *gin.Context) {
	c.JSON(200, gin.H{"message": "pong"})
}

func hanldeNotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "404 Not Found | PURGEN",
		"error": "404 Not Found",
	})
}
