package mail

import (
	"fmt"
	"github.com/go-gomail/gomail"
	"github.com/shouc/gbrm/config"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

var MAIN_HTML string
var AUTH_HTML string
var CHANNEL = make(chan *gomail.Message)

func init() {
	mainHtml, _ := ioutil.ReadFile("external/main.html")
	authHtml, _ := ioutil.ReadFile("external/authCode.html")
	MAIN_HTML = string(MAIN_HTML)
	AUTH_HTML = strings.Replace(string(mainHtml), "/@CONTENT@/",
		string(authHtml), 1)
	go func() {
		d := gomail.NewDialer(
			config.MAIL_HOST,
			80,
			config.MAIL_UN,
			config.MAIL_PW,
		)

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-CHANNEL:
				if !ok {
					return
				}
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}
					open = false
				}
			}
		}
	}()
}

func SendAuthCode(authCode int, email string) {
	authCodeStr := strconv.Itoa(authCode)
	content := strings.Replace(AUTH_HTML, "/@AUTHCODE@/", authCodeStr, 1)
	m := gomail.NewMessage()
	m.SetHeader("From", "GBRM <intl@u.zwang.tech>")
	m.SetHeader("To", fmt.Sprintf("Custormer <%s>", email))
	m.SetHeader("Subject", "Your Auth Code Arrived!")
	m.SetBody("text/html", content)
	CHANNEL <- m
}

func SendNotification(title string, r string, emails []string) {
	content := strings.Replace(MAIN_HTML, "/@CONTENT@/", r, 1)
	m := gomail.NewMessage()
	m.SetHeader("From", "GBRM <intl@u.zwang.tech>")
	m.SetHeader("To", emails...)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", content)
	CHANNEL <- m
}
