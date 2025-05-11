package regex

import (
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	"regexp"
)

func Password(pass string) error {
	if len(pass) < 8 {
		return customErr.ErrPasswordTooShort
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(pass) {
		return customErr.ErrPasswordMissingLower
	}

	if !regexp.MustCompile(`[A-Z]`).MatchString(pass) {
		return customErr.ErrPasswordMissingUpper
	}

	if !regexp.MustCompile(`\d`).MatchString(pass) {
		return customErr.ErrPasswordMissingNumber
	}

	if !regexp.MustCompile(`[@$!%*?&#]`).MatchString(pass) {
		return customErr.ErrPasswordMissingSpecial
	}

	return nil
}
