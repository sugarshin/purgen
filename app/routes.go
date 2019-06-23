package app

import (
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

	obj := gin.H{
		"url": formDataURL,
	}

	ch := make(chan http.Response)
	go makeRequest(formDataURL, ch)
	res := <-ch

	if res.Body == nil {
		obj["error"] = "Server IP address could not be found."
		c.HTML(400, "index.tmpl", obj)
		return
	}

	defer res.Body.Close()

	imageSources, err := getImageSourcesFromReader(res.Body)

	if err != nil {
		obj["error"] = "Can not found image resources."
		c.HTML(400, "index.tmpl", obj)
		return
	}

	if len(imageSources) == 0 {
		obj["message"] = "No images."
		c.HTML(200, "index.tmpl", obj)
		return
	}

	results := []Result{}

	for i := range imageSources {
		results = append(results, purge(imageSources[i], formDataURL))
	}

	obj["results"] = results
	c.HTML(200, "index.tmpl", obj)
}

func handlePing(c *gin.Context) {
	c.String(200, "pong")
}

func hanldeNotFound(c *gin.Context) {
	c.HTML(404, "index.tmpl", gin.H{
		"title":    "404 Not Found | PURGEN",
		"notfound": true,
	})
}
