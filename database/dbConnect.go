package database

import (
	"fmt"
	"github.com/shouc/gbrm/config"
	"github.com/jinzhu/gorm"

	"github.com/go-redis/redis"
	// Import MySQL database driver
	_ "github.com/jinzhu/gorm/dialects/mysql"
)


// DB global variable to access gorm
var DB *gorm.DB
var err error
var redisClient *redis.Client

// InitDB - function to initialize db
func initDB() *gorm.DB {
	var db = DB
	dbUser := config.DB_USER
	dbPassword := config.DB_PASS
	dbName := config.DB_NAME
	dbHost := config.DB_HOST

	db, err = gorm.Open(
		"mysql", dbUser+":"+ dbPassword +"@tcp(" + dbHost + ")/"+ dbName +"?charset=utf8")
	if err != nil {
		fmt.Println("DB err: ", err)
	}
	// Only for debugging
	if err == nil {
		fmt.Println("DB connection successful!")
	}
	DB = db
	return DB
}

func initRedis() *redis.Client {
	redisHost := config.REDIS_HOST
	redisPassword := config.REDIS_PASSWORD
	redisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, // no password set
		DB:       0,  // use default DB
	})
	return redisClient
}

func init(){
	initDB()
	initRedis()
}

// GetDB - get a connection
func GetDB() *gorm.DB {
	return DB
}

func GetRedis() *redis.Client {
	return redisClient
}
