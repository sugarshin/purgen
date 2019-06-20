package app

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Init : initialize app
func Init() *gin.Engine {
	app := gin.Default()
	app.LoadHTMLGlob("templates/*")
	setupRoutes(app)
	return app
}

// Run : start server
func Run() (*gin.Engine, error) {
	app := Init()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if !(port > 0) {
		port = 8000
	}

	fmt.Println("Listeninng on 0.0.0.0:" + strconv.Itoa(port))

	app.Run(":" + strconv.Itoa(port))

	return app, nil
}
