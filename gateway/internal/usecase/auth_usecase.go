package usecase

import (
	"context"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
)

type IAuthUsecase interface {
	UserRegister(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	UserLogin(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}

type AuthUsecase struct {
	client pb.AuthenticationServiceClient
}

func NewAuthUsecase(client pb.AuthenticationServiceClient) IAuthUsecase {
	return &AuthUsecase{client: client}
}

func (a *AuthUsecase) UserRegister(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
	grpcResponse, err := a.client.UserRegister(ctx, &pb.RegisterRequest{
		Name:            req.Name,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
	})
	if err != nil {
		return nil, err
	}

	return &model.RegisterResponse{
		ID:        grpcResponse.Id,
		Name:      grpcResponse.Name,
		Email:     grpcResponse.Email,
		CreatedAt: grpcResponse.CreatedAt.AsTime(),
		UpdatedAt: grpcResponse.CreatedAt.AsTime(),
	}, nil
}

func (a *AuthUsecase) UserLogin(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
	grpcResponse, err := a.client.UserLogin(ctx, &pb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		JWTToken: grpcResponse.JwtToken,
	}, nil
}
