package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS middleware
func GBRM() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		var isSuccess, _ = c.Get("success")
		var result, _ = c.Get("result")
		var message, _ = c.Get("message")
		var cache = c.GetInt("cache")
		var total = c.GetInt("total")
		if cache == 1 {
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.String(200, `{"message":null,"result":`+result.(string)+`,"success":1,"cache":1}`)
		} else {
			c.JSON(200, gin.H{
				"success": isSuccess,
				"result":  result,
				"message": message,
				"total":   total,
				"cache": 0,
			})
		}
	}
}
