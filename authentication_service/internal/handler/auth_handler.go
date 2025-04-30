package handler

import (
	"context"

	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/pb"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/usecase"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
)

type AuthHandler struct {
	pb.UnimplementedAuthenticationServiceServer
	usecase usecase.IAuthUsecase
}

func NewAuthHandler(usecase usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{
		usecase: usecase,
	}
}

func (ah *AuthHandler) UserRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	res, err := ah.usecase.UserRegister(ctx, req)
	if err != nil {
		return nil, customErr.MapErrorToStatus(err)
	}

	return res, nil
}

func (ah *AuthHandler) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	res, err := ah.usecase.UserLogin(ctx, req)
	if err != nil {
		return nil, customErr.MapErrorToStatus(err)
	}

	return res, nil
}
