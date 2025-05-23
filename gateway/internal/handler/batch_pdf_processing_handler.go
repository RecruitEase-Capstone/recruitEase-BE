package handler

import (
	"io"
	"net/http"
	"strings"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/middleware"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/response"
)

type BatchPdfProcessingHandler struct {
	usecase usecase.IBatchPdfProcessingUsecase
}

func NewBatchPdfProcessingHandler(usecase usecase.IBatchPdfProcessingUsecase) *BatchPdfProcessingHandler {
	return &BatchPdfProcessingHandler{
		usecase: usecase,
	}
}

func BatchProcessingRoutes(router *http.ServeMux,
	handler *BatchPdfProcessingHandler,
	middleware middleware.IMiddleware) {
	router.Handle("/api/cv/summarize", middleware.JwtAuthMiddleware(http.HandlerFunc(handler.HandleBatchUpload)))
	router.Handle("/api/cv/history", middleware.JwtAuthMiddleware(http.HandlerFunc(handler.FetchSummarizedPdfHistory)))
}

func (br *BatchPdfProcessingHandler) HandleBatchUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.FailedResponse(w, http.StatusMethodNotAllowed, "http method not allowed", nil)
		return
	}

	userId, err := utils.GetUserIdFromContext(w, r)

	file, fileHeader, err := r.FormFile("zipFile")
	if err != nil {
		return
	}

	if !strings.HasSuffix(strings.ToLower(fileHeader.Filename), ".zip") {
		response.FailedResponse(w, http.StatusBadRequest, "only zip file are accepted", nil)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, "file size exceeds 10 mb", nil)
		return
	}

	zipBytes, err := io.ReadAll(file)
	if err != nil {
		return
	}

	res, err := br.usecase.UnzipAndUpload(r.Context(), zipBytes, userId)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, "failed to unzip and upload batch pdf`s", err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "successfully summary batch cv", res)
}

func (br *BatchPdfProcessingHandler) FetchSummarizedPdfHistory(w http.ResponseWriter, r *http.Request) {
	userId, err := utils.GetUserIdFromContext(w, r)
	if err != nil {
		response.FailedResponse(w, http.StatusUnauthorized, "failed to get user id", err.Error())
		return
	}

	res, err := br.usecase.FetchSummarizedPdfHistory(r.Context(), userId)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, "failed to fetch summarized pdf history", err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "successfully fetch summarized pdf history", res)
}
