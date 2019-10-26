package panel

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"net/http"
)

// Views

func GetHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"title": "GBRM",
	})
}

func GetPetition(c *gin.Context) {
	cursor := database.GetDB()
	var result []model.Petitions
	cursor.Find(&result)
	jsonObj, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	c.HTML(http.StatusOK, "petition", gin.H{
		"title": "GBRM",
		"data":  string(jsonObj),
	})
}

func GetMate(c *gin.Context) {
	cursor := database.GetDB()
	var result []model.Mates
	cursor.Find(&result)
	jsonObj, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}
	c.HTML(http.StatusOK, "mate", gin.H{
		"title": "GBRM",
		"data":  string(jsonObj),
	})
}

func GetStory(c *gin.Context) {
	c.HTML(http.StatusOK, "story", gin.H{
		"title": "GBRM",
	})
}
