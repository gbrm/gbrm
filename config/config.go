package config

import (
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
)

const PROD = "prod"
const STAGING = "staging"

var DB_USER string
var DB_PASS string
var DB_NAME string
var DB_HOST string
var REDIS_HOST string
var REDIS_PASSWORD string
var SERVER_PORT string
var ALIPAY_PUB_KEY string
var ALIPAY_PRIV_KEY string
var ALI_ACCESS_ID string
var ALI_ACCESS_SECRET string
var MAIL_HOST string
var MAIL_UN string
var MAIL_PW string

func readKey() (string, string) {
	pub, _ := ioutil.ReadFile("external/alipay_public.txt")
	priv, _ := ioutil.ReadFile("external/alipay_private.txt")
	return string(pub), string(priv)
}

func init() {
	var envFile string
	if os.Getenv("stage") == PROD {
		envFile = "env/.prod_env"
	} else if os.Getenv("stage") == STAGING {
		envFile = "env/.staging_env"
	} else {
		envFile = "env/.test_env"
	}
	_ = godotenv.Load(envFile)
	DB_USER = os.Getenv("DB_USER")
	DB_PASS = os.Getenv("DB_PASS")
	DB_NAME = os.Getenv("DB_NAME")
	DB_HOST = os.Getenv("DB_HOST")

	REDIS_HOST = os.Getenv("REDIS_HOST")
	REDIS_PASSWORD = os.Getenv("REDIS_PASSWORD")
	SERVER_PORT = os.Getenv("SERVER_PORT")

	MAIL_HOST = os.Getenv("MAIL_HOST")
	MAIL_UN = os.Getenv("MAIL_UN")
	MAIL_PW = os.Getenv("MAIL_PW")
}
