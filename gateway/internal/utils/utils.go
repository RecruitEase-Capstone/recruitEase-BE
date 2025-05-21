package utils

import (
	"errors"
	"net/http"

	customContext "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/context"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/response"
)

func GetUserIdFromContext(w http.ResponseWriter, r *http.Request) (string, error) {

	userID := r.Context().Value(customContext.UserIDKey)
	if userID == "" {
		response.FailedResponse(w, http.StatusUnauthorized, "User ID tidak ditemukan dalam konteks", nil)
		return "", errors.New("user id not found in context")
	}

	stringUserID, ok := userID.(string)
	if !ok {
		response.FailedResponse(w, http.StatusBadRequest, "User ID tidak valid dalam konteks", nil)
		return "", errors.New("invalid or missing userID in context")
	}

	return stringUserID, nil
}
