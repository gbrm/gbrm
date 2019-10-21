package utils

import (
	"encoding/json"
	"fmt"
	"github.com/shouc/gbrm/config"
	"github.com/shouc/gbrm/database"
	"github.com/gin-gonic/gin"
	"os"
	"time"
)

func SetCache(c *gin.Context, data interface{}){
	if os.Getenv("stage") != config.PROD{
		return
	}
	var result []byte
	result, _ = json.Marshal(data)
	var name = c.GetString("cacheName")
	var rc = database.GetRedis()
	_ = rc.Set(name, string(result), 1 * time.Hour).Err()
}

func GetCache(c *gin.Context) bool{
	if os.Getenv("stage") == config.PROD {
		var name= c.GetString("cacheName")
		fmt.Println(name)
		var rc= database.GetRedis()
		var result, _= rc.Get(name).Result()
		if len(result) > 0 {
			c.Set("cache", 1)
			c.Set("success", 1)
			c.Set("result", result)
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

