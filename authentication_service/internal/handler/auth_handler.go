package handler

import (
	"context"

	"buf.build/go/protovalidate"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/pb/v1"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/usecase"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthenticationServiceServer
	usecase   usecase.IAuthUsecase
	validator protovalidate.Validator
}

func NewAuthHandler(usecase usecase.IAuthUsecase, validator protovalidate.Validator) *AuthHandler {
	return &AuthHandler{
		usecase:   usecase,
		validator: validator,
	}
}

func (ah *AuthHandler) UserRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := ah.validator.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	res, err := ah.usecase.UserRegister(ctx, req)
	if err != nil {
		return nil, customErr.MapErrorToStatus(err)
	}

	return res, nil
}

func (ah *AuthHandler) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := ah.validator.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	res, err := ah.usecase.UserLogin(ctx, req)
	if err != nil {
		return nil, customErr.MapErrorToStatus(err)
	}

	return res, nil
}
