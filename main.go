package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parthvinchhi/daily-data/pkg/handlers"
)

func main() {
	r := gin.Default()

	r.LoadHTMLFiles("templates/index.html")

	r.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})
	r.POST("/fetch-data", handlers.Handler)

	r.Run(":8080")
}
