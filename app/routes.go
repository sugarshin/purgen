package app

import (
	"log"
	"net/http"

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

	ch := make(chan http.Response)
	go makeRequest(formDataURL, ch)
	res := <-ch

	if res.Body == nil {
		c.HTML(400, "index.tmpl", gin.H{
			"error": "Server IP address could not be found.",
		})
		return
	}

	defer res.Body.Close()

	imageSources, err := getImageSourcesFromReader(res.Body)

	if err != nil {
		log.Fatal(err)
		c.HTML(400, "index.tmpl", gin.H{
			"error": "Can not found image resources.",
		})
		return
	}

	results := []Result{}

	for i := range imageSources {
		results = append(results, purge(imageSources[i], formDataURL))
	}

	c.HTML(200, "index.tmpl", gin.H{
		"results": results,
	})
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func hanldeNotFound(c *gin.Context) {
	c.HTML(404, "index.tmpl", gin.H{
		"title": "404 Not Found | PURGEN",
		"error": "404 Not Found",
	})
}
