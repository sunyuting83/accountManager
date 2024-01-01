package controller

import (
	BadgerDB "colaAPI/ImageServer/badger"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Path string `uri:"path" binding:"required"`
}

func GetImage(c *gin.Context) {
	var person Person
	if err := c.ShouldBindUri(&person); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":  1,
			"message": err.Error(),
		})
		return
	}
	c.Writer.Header().Add("Content-Type", "image/jpeg")

	imagePath, _ := BadgerDB.Get([]byte(person.Path))
	image, _ := os.ReadFile(string(imagePath))
	c.Writer.Write(image)
}
