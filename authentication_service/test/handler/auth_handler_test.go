package handler_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/handler"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
)

func TestUserRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockIAuthUsecase(ctrl)
	mockValidator := NewMockValidator(ctrl)

	h := handler.NewAuthHandler(mockUsecase, mockValidator)

	tests := []struct {
		name          string
		input         *pb.RegisterRequest
		setupMocks    func()
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name: "success",
			input: &pb.RegisterRequest{
				Email:    "jamal@example.com",
				Password: "Rahasia#123",
				Name:     "Jamal",
			},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(nil)

				mockUsecase.EXPECT().
					UserRegister(gomock.Any(), gomock.Any()).
					Return(&pb.RegisterResponse{Id: "uuid"}, nil)
			},
			expectedCode: codes.OK,
		},
		{
			name:  "validation error",
			input: &pb.RegisterRequest{},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(errors.New("validation failed"))
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "validation failed",
		},
		{
			name: "usecase error",
			input: &pb.RegisterRequest{
				Email:    "jamal@example.com",
				Password: "Rahasia#123",
				Name:     "Jamal",
			},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(nil)

				mockUsecase.EXPECT().
					UserRegister(gomock.Any(), gomock.Any()).
					Return(nil, customErr.ErrEmailExist)
			},
			expectedCode:  codes.AlreadyExists,
			expectedError: "email already registered",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			res, err := h.UserRegister(context.Background(), tt.input)

			if tt.expectedCode == codes.OK {
				assert.NotNil(t, res)
				assert.Nil(t, err)
			} else {
				assert.Nil(t, res)
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, s.Code())
				assert.Contains(t, s.Message(), tt.expectedError)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := NewMockIAuthUsecase(ctrl)
	mockValidator := NewMockValidator(ctrl)

	h := handler.NewAuthHandler(mockUsecase, mockValidator)

	tests := []struct {
		name          string
		input         *pb.LoginRequest
		setupMocks    func()
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name: "Success login",
			input: &pb.LoginRequest{
				Email:    "jamal@example.com",
				Password: "Rahasia#123",
			},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(nil)

				mockUsecase.EXPECT().
					UserLogin(gomock.Any(), gomock.Any()).
					Return(&pb.LoginResponse{JwtToken: "mock-token"}, nil)
			},
			expectedCode: codes.OK,
		},
		{
			name:  "Failed - validation error",
			input: &pb.LoginRequest{},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(errors.New("validation failed"))
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "validation failed",
		},
		{
			name: "Failed - invalid credentials",
			input: &pb.LoginRequest{
				Email:    "salah@example.com",
				Password: "Rahasia#123",
			},
			setupMocks: func() {
				mockValidator.EXPECT().
					Validate(gomock.Any()).
					Return(nil)

				mockUsecase.EXPECT().
					UserLogin(gomock.Any(), gomock.Any()).
					Return(nil, customErr.ErrInvalidCredentials)
			},
			expectedCode:  codes.Unauthenticated,
			expectedError: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			res, err := h.UserLogin(context.Background(), tt.input)

			if tt.expectedCode == codes.OK {
				assert.NotNil(t, res)
				assert.Nil(t, err)
			} else {
				assert.Nil(t, res)
				s, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedCode, s.Code())
				assert.Contains(t, s.Message(), tt.expectedError)
			}
		})
	}
}
