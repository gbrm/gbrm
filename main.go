package main

import (
	"fmt"
	"github.com/shouc/gbrm/config"
	"github.com/shouc/gbrm/controller/file"
	"github.com/shouc/gbrm/controller/message"
	"github.com/shouc/gbrm/controller/pay"
	"github.com/shouc/gbrm/controller/user"
	"github.com/shouc/gbrm/controller/user_admin"
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/lib/middleware"

	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"os"
)

func main() {

	fmt.Println(os.Getenv("stage"))
	db := database.GetDB()
	redis := database.GetRedis()
	defer db.Close()
	defer redis.Close()

	fmt.Println(os.Getenv("stage"))
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(middleware.CORS())
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Stage: " + os.Getenv("stage"))
	})
	v2 := router.Group("/v2/")

	v2.Use(middleware.GBRM())
	userRoute := v2.Group("user")
	{
		userRoute.POST("login", user.UserLogin)
		userRoute.GET("info", user.UserInfo)
		userRoute.POST("register", user.UserRegister)
		userRoute.POST("confirm", user.UserConfirm)
	}
	messageRouteAuth := v2.Group("message")
	messageRouteAuth.Use(middleware.Auth())
	{
		messageRouteAuth.GET("unread", message.UnreadMessages)
		messageRouteAuth.GET("read", message.ReadMessages)
	}
	_ = router.Run(":" + config.SERVER_PORT)
}
