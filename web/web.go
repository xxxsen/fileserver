package web

import (
	_ "embed"

	"github.com/gin-gonic/gin"
)

//go:embed index.html
var content string

func Index(c *gin.Context) {
	c.Writer.Write([]byte(content))
}
