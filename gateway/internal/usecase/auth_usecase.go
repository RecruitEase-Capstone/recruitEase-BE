package usecase

import (
	"context"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
)

type IAuthClient interface {
	UserRegister(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error)
	UserLogin(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error)
}

type AuthClient struct {
	client pb.AuthenticationServiceClient
}

func NewAuthClient(client pb.AuthenticationServiceClient) IAuthClient {
	return &AuthClient{client: client}
}

func (a *AuthClient) UserRegister(ctx context.Context, req *model.RegisterRequest) (*model.RegisterResponse, error) {
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

func (a *AuthClient) UserLogin(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
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
