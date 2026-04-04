package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RequestData struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

type ResponseData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	r := gin.Default()

	r.POST("/data", func(c *gin.Context) {
		var req RequestData
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Логика обработки (пример)
		resp := ResponseData{
			Status:  "ok",
			Message: "Processed user " + req.Name,
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
