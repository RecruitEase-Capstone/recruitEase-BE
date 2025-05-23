package usecase_test

import (
	context "context"
	"fmt"
	"testing"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var CTX = context.TODO()

func TestRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	ctrl.Finish()

	mockGrpc := NewMockAuthenticationServiceClient(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockGrpc)

	testCases := []struct {
		name             string
		input            *model.RegisterRequest
		mockBehavior     func(mockGrpc *MockAuthenticationServiceClient)
		expectedResponse *model.RegisterResponse
		error            error
	}{
		{
			name: "Success - Register new user",
			input: &model.RegisterRequest{
				Name:            "Jamal",
				Email:           "jamalunyu@gmail.com",
				Password:        "Rahasia#123",
				ConfirmPassword: "Rahasia#123",
			},
			mockBehavior: func(mockGrpc *MockAuthenticationServiceClient) {
				mockGrpc.EXPECT().UserRegister(CTX, gomock.Any()).
					Return(&pb.RegisterResponse{
						Id:        "uuid",
						Name:      "Jamal",
						Email:     "jamalunyu@gmail.com",
						CreatedAt: timestamppb.Now(),
						UpdatedAt: timestamppb.Now(),
					}, nil)
			},
			expectedResponse: &model.RegisterResponse{
				ID:    "uuid",
				Name:  "Jamal",
				Email: "jamalunyu@gmail.com",
			},
			error: nil,
		},
		{
			name: "Failed - Password less than 8 letters",
			input: &model.RegisterRequest{
				Name:            "Jamal",
				Email:           "jamalunyu@gmail.com",
				Password:        "kurang",
				ConfirmPassword: "kurang",
			},
			mockBehavior: func(mockGrpc *MockAuthenticationServiceClient) {
				mockGrpc.EXPECT().UserRegister(CTX, gomock.Any()).
					Return(nil, fmt.Errorf("Password must be at least 8 characters long"))
			},
			expectedResponse: nil,
			error:            fmt.Errorf("Password must be at least 8 characters long"),
		},
		{
			name: "Failed - Email already exists",
			input: &model.RegisterRequest{
				Name:            "Jamal duplicate",
				Email:           "duplicatejamal@gmail.com",
				Password:        "Rahasia#123",
				ConfirmPassword: "Rahasia#123",
			},
			mockBehavior: func(mockGrpc *MockAuthenticationServiceClient) {
				mockGrpc.EXPECT().
					UserRegister(CTX, gomock.Any()).
					Return(nil, fmt.Errorf("email already exist"))
			},
			expectedResponse: nil,
			error:            fmt.Errorf("email already exist"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mockGrpc)

			res, err := authUsecase.UserRegister(CTX, tc.input)
			if tc.error != nil {
				assert.EqualError(t, err, tc.error.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.expectedResponse.ID, res.ID)
				assert.Equal(t, tc.expectedResponse.Name, res.Name)
				assert.Equal(t, tc.expectedResponse.Email, res.Email)
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpc := NewMockAuthenticationServiceClient(ctrl)
	authUsecase := usecase.NewAuthUsecase(mockGrpc)

	testCases := []struct {
		name             string
		input            *model.LoginRequest
		mockBehavior     func(mockGrpc *MockAuthenticationServiceClient)
		expectedResponse *model.LoginResponse
		error            error
	}{
		{
			name: "Sucess - Login with registered account",
			input: &model.LoginRequest{
				Email:    "jamalunyu@gmail.com",
				Password: "Rahasia#123",
			},
			mockBehavior: func(mockGrpc *MockAuthenticationServiceClient) {
				mockGrpc.EXPECT().
					UserLogin(CTX, gomock.Any()).
					Return(&pb.LoginResponse{
						JwtToken: "jwt-token",
					}, nil)
			},
			expectedResponse: &model.LoginResponse{
				JWTToken: "jwt-token",
			},
			error: nil,
		},
		{
			name: "Failed - Email not found",
			input: &model.LoginRequest{
				Email:    "notfound@gmail.com",
				Password: "Rahasia#123",
			},
			mockBehavior: func(mockGrpc *MockAuthenticationServiceClient) {
				mockGrpc.EXPECT().
					UserLogin(CTX, gomock.Any()).
					Return(nil, fmt.Errorf("invalid email or password"))
			},
			error: fmt.Errorf("invalid email or password"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mockGrpc)

			res, err := authUsecase.UserLogin(CTX, tc.input)
			if tc.error != nil {
				assert.EqualError(t, err, tc.error.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.JWTToken)
			}
		})
	}
}
