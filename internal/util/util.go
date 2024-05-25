package util

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Response is customized to help return all responses need
func Response(c *gin.Context, message string, status int, data interface{}, errs []string) {
	responsedata := gin.H{
		"message":   message,
		"data":      data,
		"errors":    errs,
		"status":    http.StatusText(status),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.IndentedJSON(status, responsedata)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
