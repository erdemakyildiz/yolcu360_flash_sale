package utils

import (
	"fmt"
	"github.com/gofiber/fiber/v2/log"
	"time"
)

func CreateLogMessage(msg string, err error) string {
	msg = fmt.Sprintf("[%s] %s : %v", time.Now().String(), msg, err)
	log.Error(msg)
	return msg
}
