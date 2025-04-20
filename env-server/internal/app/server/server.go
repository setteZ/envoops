package server

import (
	"env-server/utils"

	"github.com/gin-gonic/gin"
)

var DEFAULT_PORT = "8080"

func Run() {
	r := setupRouter(gin.DebugMode)
	r.Run(":" + utils.GetEnv("ENVSERVER_PORT", DEFAULT_PORT))
}
