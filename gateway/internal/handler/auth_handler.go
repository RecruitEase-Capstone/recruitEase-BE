package handler

import (
	"encoding/json"
	"net/http"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/error"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/response"
)

type AuthHandler struct {
	usecase usecase.IAuthUsecase
}

func NewAuthHandler(usecase usecase.IAuthUsecase) *AuthHandler {
	return &AuthHandler{usecase: usecase}
}

func AuthRoutes(router *http.ServeMux, authHandler *AuthHandler) {
	router.Handle("/api/auth/register", http.HandlerFunc(authHandler.UserRegister))
	router.Handle("/api/auth/login", http.HandlerFunc(authHandler.UserLogin))
}

func (ah *AuthHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	var req *model.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	res, err := ah.usecase.UserRegister(r.Context(), req)
	if err != nil {
		statusCode, msg := customErr.GRPCErrorToHTTP(err)
		response.FailedResponse(w, statusCode, msg, nil)
		return
	}

	response.SuccessResponse(w, http.StatusCreated, "successfully create new account", res)
}

func (ah *AuthHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var req *model.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error(), nil)
		return
	}

	res, err := ah.usecase.UserLogin(r.Context(), req)
	if err != nil {
		statusCode, msg := customErr.GRPCErrorToHTTP(err)
		response.FailedResponse(w, statusCode, msg, nil)
		return
	}

	response.SuccessResponse(w, http.StatusOK, "successfully login to account", res)
}
