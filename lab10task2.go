package main

import (
    "log"
    "time"
    "github.com/gin-gonic/gin"
)

//Middleware
func loggerMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        method := c.Request.Method

        c.Next()

        status := c.Writer.Status()
        duration := time.Since(start)

        log.Printf("[%s] %s - status: %d - duration: %v", method, path, status, duration)
    }
}

func main() {
    r := gin.Default()

    r.Use(loggerMiddleware())

    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    r.GET("/hello/:name", func(c *gin.Context) {
        name := c.Param("name")
        c.JSON(200, gin.H{"hello": name})
    })

    r.Run(":8080")
}