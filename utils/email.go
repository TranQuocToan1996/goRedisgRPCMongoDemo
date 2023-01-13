package utils

import (
	"net/mail"
	"regexp"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
)

func IsEmail(email string) bool {
	if emailRegex == nil {
		_, err := mail.ParseAddress(email)
		return err == nil
	}
	return emailRegex.MatchString(email)
}
