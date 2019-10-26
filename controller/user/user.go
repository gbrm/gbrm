package user

import (
	"github.com/gin-gonic/gin"
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"github.com/shouc/gbrm/mail"
	"github.com/shouc/gbrm/utils"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// Views

func GetLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login", gin.H{
		"title": "Login | GBRM",
	})
}

func PostLogin(c *gin.Context) {
	var userParam string
	userParam = c.PostForm("email")
	var password string
	password = utils.Encrypt(c.PostForm("password"))
	var cursor = database.GetDB()
	var result struct {
		Session string `json:"session"`
		Id      int    `json:"user_id"`
		Type    int    `json:"type"`
		Token   string `json:"token"`
	}
	cursor.
		Table("users").
		Where("password = ?", password).
		Or("email = ?", userParam).
		Limit("1").
		Scan(&result)
	if result.Id == 0 {
		c.HTML(http.StatusOK, "login", gin.H{
			"title":   "Login | GBRM",
			"message": "Cannot find user!",
		})
	} else {
		result.Session = utils.AddSession(result.Token)
		c.Redirect(http.StatusFound, "home?session="+result.Session)
	}
}

func GetRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "register", gin.H{
		"title": "Register | GBRM",
	})
}

func PostRegister(c *gin.Context) {
	var email string
	email = c.PostForm("email")
	if c.PostForm("password") != c.PostForm("confirm_password") {
		c.HTML(http.StatusOK, "register", gin.H{
			"title":   "Register | GBRM",
			"message": "Password Not Match",
			"email":   email,
		})
	}
	var password string
	password = utils.Encrypt(c.PostForm("password"))

	// Insert Section
	user := model.Users{
		Email:    email,
		Password: password,
		IsAuth:   false,
		Type:     1,
		ProTill:  0}
	var cursor = database.GetDB()
	insertErr := cursor.Create(&user)
	if insertErr.RowsAffected == 1 {
		addAuthCode(email)
		c.Redirect(http.StatusFound, "confirm?email="+user.Email)
	} else {
		errMessage := insertErr.Error.Error()
		var message string
		if strings.Contains(errMessage, "key 'email'") {
			message = "Email has already been taken!"
		} else {
			message = "System Error"
		}
		c.HTML(http.StatusOK, "register", gin.H{
			"title":   "Register | GBRM",
			"message": message,
			"email":   email,
		})
	}
}

func GetConfirm(c *gin.Context) {
	userMail := c.Query("email")
	c.HTML(http.StatusOK, "confirm", gin.H{
		"title": "Confirm Your Email | GBRM",
		"email": userMail,
	})
}

func PostConfirm(c *gin.Context) {
	var authCode = c.DefaultPostForm("auth_code", "0")
	var email = c.DefaultPostForm("email", "0")

	var rc = database.GetRedis()
	authCodeSaved, _ := rc.Get(email).Result()
	if authCodeSaved == "" {
		c.HTML(http.StatusOK, "confirm", gin.H{
			"title":   "Confirm | GBRM",
			"email":   email,
			"message": "The auth code is already expired or already used",
		})
	} else {
		if authCodeSaved == authCode {
			var token, session = utils.AddSessionWithoutToken()
			var user model.Users
			cursor := database.GetDB()
			cursor.Model(&user).
				Where("email = ?", email).
				Updates(model.Users{Token: token, IsAuth: false})
			_, _ = rc.Del(email).Result()
			c.Redirect(http.StatusFound, "home?session="+session)
		} else {
			c.HTML(http.StatusOK, "confirm", gin.H{
				"title":   "Confirm | GBRM",
				"email":   email,
				"message": "Incorrect Auth Code",
			})
		}
	}
}

// API

func UserInfo(c *gin.Context) {
	var session string
	session = c.Query("session")
	type resultStruct struct {
		Type   int  `json:"type"`
		IsPaid bool `json:"is_paid"`
	}
	var user = utils.RetrieveSession(session)
	if user.Id == 0 {
		c.Set("success", 0)
		c.Set("message", "Cannot find user")
	} else {
		c.Set("success", 1)
		c.Set("result", resultStruct{
			Type: user.Type,
		})
	}
}

func addAuthCode(email string) string {
	var rc = database.GetRedis()
	var authCode = rand.Intn(100000)
	mail.SendAuthCode(authCode, email)
	_ = rc.Set(email, authCode, time.Hour).Err()
	return string(authCode)
}
