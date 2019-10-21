package middleware

import (
	"github.com/shouc/gbrm/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func getRequestRoute(c *gin.Context) (string, string) {
	urls := strings.Split(c.Request.URL.Path, "/")
	if len(urls) >= 3 {
		return urls[2], urls[3]
	}
	if len(urls) >= 2 {
		return urls[2], ""
	}
	return "", ""
}

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		var session string
		switch c.Request.Method {
		case "GET":
			session = c.DefaultQuery("session", "")
			break
		case "DELETE":
			session = c.DefaultQuery("session", "")
			break
		default:
			session = c.DefaultPostForm("session", "")
		}
		userObj := utils.RetrieveSession(session)
		baseRoute, subRoute := getRequestRoute(c)
		if userObj.Id == 0 {
			c.Set("message", "Your login has expired")
			c.Set("success", 0)
			c.Set("isAuth", false)
			// admin routes
		} else if userObj.Type != 5 && baseRoute == "admin" {
			c.Set("message", "Your are not admin")
			c.Set("success", 0)
			c.Set("isAuth", false)
			// contrib routes
		} else if userObj.Type < 3 && baseRoute == "contributor" {
			if userObj.Type == 3 && subRoute == "video" {
				c.Set("message", "Your are normal contributor")
				c.Set("success", 0)
				c.Set("isAuth", false)
				return
			}
			c.Set("message", "Your are not contributor")
			c.Set("success", 0)
			c.Set("isAuth", false)
			return
		} else {
			c.Set("userObj", &userObj)
			c.Set("isAuth", true)
			c.Set("userId", userObj.Id)
			return
		}
		c.Next()

	}
}
