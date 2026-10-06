// Задание 6. Go (Gin) — один из двух сервисов для сравнения скорости под нагрузкой.
//
// Служебный эндпоинт без логирования и бизнес-логики, чтобы замерялась
// чистая стоимость HTTP-обработки.
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	log.Println("Gin-сервис задания 6 слушает http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("сервер не запустился: %v", err)
	}
}
