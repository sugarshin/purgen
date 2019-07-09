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

	// req, err := http.NewRequest("GET", formDataURL, nil)

	// if err != nil {
	// 	obj["error"] = "InternalServerError"
	// 	c.HTML(400, "index.tmpl", obj)
	// 	return
	// }

	// client := &http.Client{
	// 	Timeout: time.Duration(10 * time.Second),
	// }

	ch := make(chan http.Response)
	go makeGetRequest(formDataURL, ch)
	res := <-ch

	if res.Body == nil {
		obj["error"] = "Server IP address could not be found."
		c.HTML(400, "index.tmpl", obj)
		return
	}

	defer res.Body.Close()

	resources, err := getResourcesFromReader(res.Body)

	if err != nil {
		obj["error"] = "Can not found resources."
		c.HTML(400, "index.tmpl", obj)
		return
	}

	if len(resources) == 0 {
		obj["message"] = "No resources."
		c.HTML(200, "index.tmpl", obj)
		return
	}

	results := []Result{}

	for i := range resources {
		results = append(results, purge(resources[i], formDataURL))
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
