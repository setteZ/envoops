package server

import (
	"env-server/internal/handlers"
	"env-server/internal/version"

	"github.com/gin-gonic/gin"
)

func setupRouter(relMode string) *gin.Engine {
	gin.SetMode(relMode)
	router := gin.Default()

	/* frontend route */
	router.GET("/", handlers.HomePage)
	/******************/

	/*** _R__ api ***/
	router.GET("/api/:table", handlers.ReadTable)
	/****************/

	/* test route */
	router.GET("/ping", handlers.PingResponse)
	/****************/

	/* server info route */
	router.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"version": version.Version,
			"commit":  version.Commit,
			"date":    version.Date,
		})
	})
	/****************/

	/* static route */
	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	/****************/
	return router
}
