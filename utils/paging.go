package utils

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetPageFromContext(c *gin.Context) (int, int) {
	page := c.DefaultQuery("page", "1")
	pageInt, _ := strconv.Atoi(page)
	limit := 5
	offset := (pageInt - 1) * limit
	return offset, limit
}
