package index

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Views

func GetIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "GBRM",
	})
}
