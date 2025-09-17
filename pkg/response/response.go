package response

import "github.com/gin-gonic/gin"

func OK(c *gin.Context, data any) {
	c.JSON(200, data)
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": gin.H{
			"code":    code,
			"message": message,
		},
	})
}
