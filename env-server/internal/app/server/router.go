package server

import (
	"env-server/internal/handlers"

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

	/* static route */
	router.StaticFile("/favicon.ico", "./web/assets/favicon.ico")
	/****************/
	return router
}
