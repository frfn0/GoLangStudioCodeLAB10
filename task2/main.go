// Задание 2. Middleware для логирования в Go.
//
// Сервис поднимает два эндпоинта и пишет в лог каждый запрос через
// собственный middleware: метод, путь, код ответа и время обработки.
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware логирует каждый обработанный запрос.
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path

		c.Next()

		log.Printf(
			"[GIN] %s %s -> %d (%s)",
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start),
		)
	}
}

func main() {
	// gin.New() без встроенного логгера: логирование выполняет наш middleware,
	// иначе каждый запрос писался бы в stdout дважды.
	r := gin.New()
	r.Use(LoggerMiddleware(), gin.Recovery())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": c.Param("name")})
	})

	log.Println("Gin-сервис слушает http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("сервер не запустился: %v", err)
	}
}
