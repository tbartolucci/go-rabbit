package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Addition struct {
	Number1 int64 `json:"number1" binding:"required"`
	Number2 int64 `json:"number2" binding:"required"`
}


func main() {
	r := gin.Default()
	r.GET("/hi", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Josh",
		})
	})
	r.POST("/add", func(c *gin.Context){
		var json Addition
		if err := c.ShouldBindJSON(&json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}


	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
