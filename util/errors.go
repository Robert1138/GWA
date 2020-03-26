package util

import (
	"errors"
)

var (
	// ErrUserExists returns a user already exists error which occurs when registering a user whose username is alreadu in use
	ErrUserExists = errors.New("user already exists")
	// ErrPasswordInvalid occurs when supplied password in User struct is less than 8
	ErrPasswordInvalidFormat = errors.New("password is too short")
)
