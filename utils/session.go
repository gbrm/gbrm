package utils

import (
	"crypto/sha256"
	"fmt"
	"github.com/shouc/gbrm/database"
	"github.com/shouc/gbrm/database/model"
	"github.com/satori/go.uuid"
	"time"
)

func Encrypt(content string) string{
	encryptor := sha256.New()
	encryptor.Write([]byte(content))
	hash := encryptor.Sum(nil)
	return fmt.Sprintf("%x", hash)
}

func GenerateUUID() string{
	val := uuid.NewV4()
	return val.String()
}

func AddSession(token string) string{
	var session = GenerateUUID()
	var rc = database.GetRedis()
	err := rc.Set(session, token, 24 * 7 * time.Hour).Err()
	fmt.Println(err)
	return session
}

func AddSessionWithoutToken() (string, string){
	var token = GenerateUUID()
	var session = GenerateUUID()
	var rc = database.GetRedis()
	err := rc.Set(session, token, 24 * 7 * time.Hour).Err()
	fmt.Println(err)
	return token, session
}

func RetrieveSession(session string) model.Users {
	var rc = database.GetRedis()
	token, _ := rc.Get(session).Result()
	cursor := database.GetDB()
	var user = model.Users{}
	cursor.
		Where("token = ?", token).
		First(&user)
	return user
}

func GetUserIdFromSession(session string) int {
	var userObj model.Users
	userObj = RetrieveSession(session)
	return userObj.Id
}
