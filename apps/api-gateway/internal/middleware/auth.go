package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")
		if err != nil || token == "" {
			c.Next()
			return
		}

		client := &http.Client{}

		req, err := http.NewRequest("GET", "http://auth:8080/refresh", nil)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to create auth request",
			})
			return
		}

		req.AddCookie(&http.Cookie{Name: "token", Value: token})

		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"error": "Auth service unavailable",
			})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		if newToken := resp.Header.Get("X-Access-Token"); newToken != "" {
			c.Set("X-Access-Token", newToken)
		}

		c.Next()
	}
}
