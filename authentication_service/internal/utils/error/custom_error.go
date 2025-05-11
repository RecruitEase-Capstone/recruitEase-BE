package error

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailExist             = errors.New("email already exist")
	ErrNotVerified            = errors.New("account has not been verified")
	ErrIncorrectPassword      = errors.New("incorrect password")
	ErrDatabase               = errors.New("database error")
	ErrRowsAffected           = errors.New("error due to there is no or more than 1 affected column")
	ErrPasswordTooShort       = errors.New("password must be at least 8 characters long")
	ErrPasswordMissingLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordMissingUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordMissingNumber  = errors.New("password must contain at least one number")
	ErrPasswordMissingSpecial = errors.New("password must contain at least one special character")
)

func MapErrorToStatus(err error) error {
	switch {
	case errors.Is(err, ErrEmailExist):
		return status.Error(codes.AlreadyExists, "email already registered")
	case errors.Is(err, ErrUserNotFound):
		return status.Error(codes.NotFound, "user not found")
	case errors.Is(err, ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid email or password")
	case errors.Is(err, ErrPasswordTooShort),
		errors.Is(err, ErrPasswordMissingLower),
		errors.Is(err, ErrPasswordMissingUpper),
		errors.Is(err, ErrPasswordMissingNumber),
		errors.Is(err, ErrPasswordMissingSpecial):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal server error")
	}
}
