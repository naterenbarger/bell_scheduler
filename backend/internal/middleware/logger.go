package middleware

import (
    "fmt"
    "time"

    "github.com/gin-gonic/gin"
)

// Logger creates a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Start timer
        start := time.Now()

        // Process request
        c.Next()

        // Log details
        duration := time.Since(start)
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path
        statusCode := c.Writer.Status()

        // Get username if available
        username := "anonymous"
        if user, exists := c.Get("username"); exists {
            username = user.(string)
        }

        // Log the request
        fmt.Printf("[%s] %s %s %d %s %s %s\n",
            time.Now().Format("2006-01-02 15:04:05"),
            method,
            path,
            statusCode,
            duration,
            clientIP,
            username,
        )
    }
} 