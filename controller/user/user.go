package user

import (
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"github.com/shouc/gbrm/mail"
	"github.com/shouc/gbrm/utils"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strings"
	"time"
)

func UserLogin(c *gin.Context) {
	var userParam string
	userParam = c.PostForm("user_param")
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
		Where("name = ?", userParam).
		Or("email = ?", userParam).
		Limit("1").
		Scan(&result)
	if result.Id == 0 {
		c.Set("success", 0)
		c.Set("message", "Cannot find user!")
	} else {
		result.Session = utils.AddSession(result.Token)
		c.Set("success", 1)
		c.Set("result", result)
	}
}

func UserInfo(c *gin.Context) {
	var session string
	session = c.Query("session")
	type resultStruct struct{
		Type int `json:"type"`
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

func UserRegister(c *gin.Context) {
	var email string
	email = c.PostForm("email")
	var username string
	username = c.PostForm("username")
	var password string
	password = utils.Encrypt(c.PostForm("password"))

	//Check Section

	// Insert Section
	user := model.Users{
		Name:     username,
		Email:    email,
		Password: password,
		IsAuth:   false,
		Type:     1,
		ProTill:  0}
	var cursor = database.GetDB()
	insertErr := cursor.Create(&user)
	if insertErr.RowsAffected == 1 {
		addAuthCode(email)
		c.Set("success", 1)
		c.Set("message", "Check your email")
	} else {
		errMessage := insertErr.Error.Error()
		c.Set("success", 0)
		if strings.Contains(errMessage, "key 'name'") {
			c.Set("message", "Username has already been taken!")
		} else if strings.Contains(errMessage, "key 'email'") {
			c.Set("message", "Email has already been taken!")
		} else {
			c.Set("message", "System Error")
		}
	}
}

func UserConfirm(c *gin.Context) {
	var authCode = c.DefaultPostForm("auth_code", "0")
	var email = c.DefaultPostForm("email", "0")

	var rc = database.GetRedis()
	authCodeSaved, _ := rc.Get(email).Result()
	if authCodeSaved == "" {
		c.Set("success", 0)
		c.Set("message", "The auth code is already expired or already used")
	} else {
		if authCodeSaved == authCode {
			var token, session = utils.AddSessionWithoutToken()
			var user model.Users
			cursor := database.GetDB()
			cursor.Model(&user).
				Where("email = ?", email).
				Updates(model.Users{Token: token, IsAuth: false})
			_, _ = rc.Del(email).Result()
			c.Set("success", 1)
			c.Set("result", session)
		} else {
			c.Set("success", 0)
			c.Set("message", "Incorrect Auth Code")
		}
	}
}
