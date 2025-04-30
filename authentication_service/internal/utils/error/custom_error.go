package error

import (
	"errors"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/regex"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExist         = errors.New("email already exist")
	ErrNotVerified        = errors.New("account has not been verified")
	ErrIncorrectPassword  = errors.New("incorrect password")
	ErrDatabase           = errors.New("database error")
	ErrRowsAffected       = errors.New("error due to there is no or more than 1 affected column")
)

func MapErrorToStatus(err error) error {
	switch {
	case errors.Is(err, ErrEmailExist):
		return status.Error(codes.AlreadyExists, "email already registered")
	case errors.Is(err, ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid email or password")
	// mapping regex errors
	case errors.Is(err, regex.ErrPasswordTooShort),
		errors.Is(err, regex.ErrPasswordMissingLower),
		errors.Is(err, regex.ErrPasswordMissingUpper),
		errors.Is(err, regex.ErrPasswordMissingNumber),
		errors.Is(err, regex.ErrPasswordMissingSpecial):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
