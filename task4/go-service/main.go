// Задание 4. Go-сервис, к которому обращается Python-сервис по HTTP.
//
// Принимает JSON на POST /data, выполняет обработку и возвращает JSON.
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	UserID int    `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

type UserResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/data", func(c *gin.Context) {
		var req UserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.Age < 0 || req.Age > 150 {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "age must be in range 0..150"})
			return
		}

		c.JSON(http.StatusOK, UserResponse{
			Status:  "ok",
			Message: "Processed user " + req.Name,
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	log.Println("Go-сервис задания 4 слушает http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("сервер не запустился: %v", err)
	}
}
