// regex/password.go
package regex

import (
	"errors"
	"regexp"
)

var (
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters long")
	ErrPasswordMissingLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingNumber  = errors.New("password must contain at least one number")
	ErrPasswordMissingSpecial = errors.New("password must contain at least one special character")
)

func Password(pass string) error {
	if len(pass) < 8 {
		return ErrPasswordTooShort
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
		return ErrPasswordMissingLower
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
		return ErrPasswordMissingUpper
	}

	if !regexp.MustCompile(`\d`).MatchString(pass) {
		return ErrPasswordMissingNumber
	}

	if !regexp.MustCompile(`[@$!%*?&#]`).MatchString(pass) {
		return ErrPasswordMissingSpecial
	}

	return nil
}
