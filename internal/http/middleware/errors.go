package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		// if handlers wrote an error, format it
		if len(c.Errors) > 0 {
			c.JSON(-1, gin.H{
				"error": gin.H{
					"code":    "INTERNAL",
					"message": c.Errors[0].Error(),
				},
			})
		}
	}
}

func JSONError(c *gin.Context, code int, msg string) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": gin.H{
			"code":    http.StatusText(code),
			"message": msg,
		},
	})
}
