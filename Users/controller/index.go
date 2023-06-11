package controller

import (
	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	// c.HTML(http.StatusOK, "index.html", nil)
	c.String(200, "ok")
}
