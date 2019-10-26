package main

import (
	"fmt"
	"github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"github.com/shouc/gbrm/config"
	"github.com/shouc/gbrm/controller/index"
	"github.com/shouc/gbrm/controller/panel"
	"github.com/shouc/gbrm/controller/user"
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"github.com/shouc/gbrm/lib/middleware"
	"os"
)

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	layouts := []string{
		"views/layout/header.html",
		"views/layout/footer.html",
	}
	r.AddFromFiles("index", "views/index.html", layouts[0], layouts[1])
	r.AddFromFiles("login", "views/user/user_base.html", "views/user/login.html", layouts[0], layouts[1])
	r.AddFromFiles("register", "views/user/user_base.html", "views/user/register.html", layouts[0], layouts[1])
	r.AddFromFiles("confirm", "views/user/user_base.html", "views/user/confirm.html", layouts[0], layouts[1])

	r.AddFromFiles("home", "views/panel/panel_base.html", "views/panel/home.html", layouts[0], layouts[1])
	r.AddFromFiles("mate", "views/panel/panel_base.html", "views/panel/mate.html", layouts[0], layouts[1])
	r.AddFromFiles("story", "views/panel/panel_base.html", "views/panel/story.html", layouts[0], layouts[1])
	r.AddFromFiles("petition", "views/panel/panel_base.html", "views/panel/petition.html", layouts[0], layouts[1])

	return r
}

func main() {

	fmt.Println(os.Getenv("stage"))
	db := database.GetDB()
	redis := database.GetRedis()
	defer db.Close()
	defer redis.Close()
	db.CreateTable(&model.Mates{})
	db.CreateTable(&model.Users{})
	db.CreateTable(&model.Messages{})
	db.CreateTable(&model.Petitions{})
	db.CreateTable(&model.Tags{})

	fmt.Println(os.Getenv("stage"))
	router := gin.Default()
	router.Use(sentrygin.New(sentrygin.Options{}))
	router.Use(middleware.CORS())

	v2 := router.Group("/v2/")

	v2.Use(middleware.GBRM())
	userRoute := v2.Group("user")
	{
		userRoute.GET("info", user.UserInfo)
	}

	router.HTMLRender = createMyRender()
	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	router.GET("login", user.GetLogin)
	router.GET("register", user.GetRegister)
	router.GET("confirm", user.GetConfirm)

	router.POST("login", user.PostLogin)
	router.POST("register", user.PostRegister)
	router.POST("confirm", user.PostConfirm)

	router.GET("", index.GetIndex)

	router.GET("home", panel.GetHome)
	router.GET("mate", panel.GetMate)
	router.GET("story", panel.GetStory)
	router.GET("petition", panel.GetPetition)

	_ = router.Run(":" + config.SERVER_PORT)
}
