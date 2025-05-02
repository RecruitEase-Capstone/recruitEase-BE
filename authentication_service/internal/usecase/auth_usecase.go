package usecase

import (
	"context"
	"time"

	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/model"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/pb/v1"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/repository"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/jwt"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/regex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IAuthUsecase interface {
	UserRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error)
	UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error)
}

type AuthUsecase struct {
	repo repository.IAuthRepository
	jwt  jwt.JWTItf
}

func NewAuthUsecase(repo repository.IAuthRepository, jwt jwt.JWTItf) IAuthUsecase {
	return &AuthUsecase{repo: repo, jwt: jwt}
}

func (au *AuthUsecase) UserRegister(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	existingUser, err := au.repo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, customErr.ErrEmailExist
	}

	if err := regex.Password(req.Password); err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createdUser := &model.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := au.repo.CreateUser(ctx, createdUser); err != nil {
		return nil, err
	}

	response := &pb.RegisterResponse{
		Id:        createdUser.ID,
		Name:      createdUser.Name,
		Email:     createdUser.Email,
		CreatedAt: timestamppb.New(createdUser.CreatedAt),
		UpdatedAt: timestamppb.New(createdUser.UpdatedAt),
	}

	return response, nil
}

func (au *AuthUsecase) UserLogin(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := au.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, customErr.ErrInvalidCredentials
	}

	token, err := au.jwt.CreateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		JwtToken: token,
	}, nil
}
