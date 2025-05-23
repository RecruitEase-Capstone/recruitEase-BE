package usecase_test

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/error"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/minio/minio-go/v7"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	bucketName string = "test-bucket"
	userId     string = "user-uuid"
	cvCtx             = context.TODO()
)

func setupZipInput() []byte {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)
	w, _ := zipWriter.Create("dummy-cv.pdf")
	w.Write([]byte("fake-pdf"))
	zipWriter.Close()
	return buf.Bytes()
}

func TestUnzipAndUpload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockGrpc := NewMockCVProcessorServiceClient(ctrl)
	mockMinio := NewMockIMinio(ctrl)
	cvUsecase := usecase.NewBatchPdfProcessing(mockMinio, mockGrpc, bucketName)

	testCases := []struct {
		name          string
		inputZipBytes []byte
		inputUserId   string
		mockBehavior  func(mockGrpc *MockCVProcessorServiceClient,
			mockMinio *MockIMinio)
		expectedResponse []*model.CVSummarizeResponse
		error            error
	}{
		{
			name:          "Success Unzip and Upload batch pdf files",
			inputZipBytes: setupZipInput(),
			inputUserId:   userId,
			mockBehavior: func(mockGrpc *MockCVProcessorServiceClient, mockMinio *MockIMinio) {
				mockMinio.EXPECT().
					MakeBucket(cvCtx, gomock.Any()).
					Return(nil)

				mockMinio.
					EXPECT().
					UploadPDF(cvCtx, bucketName, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(minio.UploadInfo{Size: 2048}, nil)

				mockGrpc.EXPECT().
					ProcessBatchPDF(cvCtx, gomock.Any()).
					Return(&pb.BatchPDFProcessResponse{
						BatchId:    gomock.Nil().String(),
						TotalFiles: 1,
						Predictions: []*pb.PredictionResult{
							{
								FileName: gomock.Any().String(),
								Prediction: &pb.CVPrediction{
									Name:              []string{"Jamal"},
									CollegeName:       []string{"UB"},
									Degree:            []string{"Computer Science"},
									GraduationYear:    []string{"2024"},
									YearsOfExperience: []string{"2 years"},
									CompaniesWorkedAt: []string{"Google"},
									Designation:       []string{"DevOps"},
									Skills:            []string{"Rust"},
									Location:          []string{"Depok"},
									EmailAddress:      []string{"Jamalunyu@gmail.com"},
								},
							},
						},
					}, nil)
			},
			expectedResponse: []*model.CVSummarizeResponse{
				{
					Name:              []string{"Jamal"},
					CollegeName:       []string{"UB"},
					Degree:            []string{"Computer Science"},
					GraduationYear:    []string{"2024"},
					YearsOfExperience: []string{"2 years"},
					CompaniesWorkedAt: []string{"Google"},
					Designation:       []string{"DevOps"},
					Skills:            []string{"Rust"},
					Location:          []string{"Depok"},
					EmailAddress:      []string{"Jamalunyu@gmail.com"},
				},
			},
			error: nil,
		},
		{
			name:          "Failed - Minio upload error",
			inputZipBytes: setupZipInput(),
			inputUserId:   userId,
			mockBehavior: func(mockGrpc *MockCVProcessorServiceClient, mockMinio *MockIMinio) {
				mockMinio.EXPECT().
					MakeBucket(cvCtx, bucketName).
					Return(nil)

				mockMinio.EXPECT().
					UploadPDF(cvCtx, bucketName, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(minio.UploadInfo{}, customErr.ErrFailedToUpload)
			},
			expectedResponse: nil,
			error:            customErr.ErrFailedToUpload,
		},
		{
			name:          "Failed - error occurs when processing pdf from grpc",
			inputZipBytes: setupZipInput(),
			inputUserId:   userId,
			mockBehavior: func(mockGrpc *MockCVProcessorServiceClient, mockMinio *MockIMinio) {
				mockMinio.EXPECT().
					MakeBucket(cvCtx, gomock.Any()).
					Return(nil)

				mockMinio.
					EXPECT().
					UploadPDF(cvCtx, bucketName, gomock.Any(), gomock.Any(), gomock.Any()).
					Return(minio.UploadInfo{Size: 2048}, nil)

				mockGrpc.EXPECT().
					ProcessBatchPDF(cvCtx, gomock.Any()).
					Return(nil, fmt.Errorf("failed to process batch PDF via gRPC"))
			},
			expectedResponse: nil,
			error:            fmt.Errorf("failed to process batch PDF via gRPC"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockBehavior(mockGrpc, mockMinio)

			res, err := cvUsecase.UnzipAndUpload(CTX, tc.inputZipBytes, "user-uuid")
			if tc.error != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.error.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, res, tc.expectedResponse)
			}
		})
	}
}
