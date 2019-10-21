package utils

import (
	"time"
)

func CurrentTime() int{
	return int(time.Now().Unix())
}
