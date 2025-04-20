package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingResponse(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"msg": "pong"})
}

func HomePage(c *gin.Context) {

	output := `<!DOCTYPE html>
<html>
<head>
    <title>env-ops</title>
</head>
<body>
    Server up
</body>
</html>
`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(output))
}

func ReadTable(c *gin.Context) {
	table := c.Param("table")

	switch table {
	case "test":
		c.JSON(http.StatusOK, gin.H{"table": table, "content": []map[string]any{gin.H{"id": 0, "key1": "value1-0", "key2": "value2-0"}, gin.H{"id": 1, "key1": "value1-1", "key2": "value2-1"}}})
	default:
		c.JSON(http.StatusNotFound, gin.H{"table": table, "msg": "missing table"})
	}
}
