package utils

import (
	"github.com/gin-gonic/gin"
)

func GetUserIdFromContext(c *gin.Context) int{
	userIdStr, _ := c.Get("userId")
	userId := int(userIdStr.(int))
	return userId
}
