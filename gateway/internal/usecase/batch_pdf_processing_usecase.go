package usecase

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/model"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/minio"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IBatchPdfProcessingUsecase interface {
	UnzipAndUpload(ctx context.Context, zipBytes []byte, userId string) (*model.CVSummarizeResponse, error)
	FetchSummarizedPdfHistory(ctx context.Context, userId string) (*model.CVSummarizeResponse, error)
}

type BatchPdfProcessingUsecase struct {
	minio      minio.IMinio
	client     pb.CVProcessorServiceClient
	bucketName string
}

func NewBatchPdfProcessing(minio minio.IMinio,
	client pb.CVProcessorServiceClient,
	bucketName string) IBatchPdfProcessingUsecase {
	return &BatchPdfProcessingUsecase{minio: minio, client: client, bucketName: bucketName}
}

func (bu *BatchPdfProcessingUsecase) UnzipAndUpload(ctx context.Context, zipBytes []byte, userId string) (*model.CVSummarizeResponse, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}

	var UploadedFiles []*pb.PDFFileInfo

	for _, zipFile := range zipReader.File {
		if zipFile.FileInfo().IsDir() {
			continue
		}

		if !strings.HasSuffix(strings.ToLower(zipFile.Name), ".pdf") {
			continue
		}

		res, err := bu.processPDF(ctx, zipFile)
		if err != nil {
			return nil, err
		}

		UploadedFiles = append(UploadedFiles, res)
	}
	batchId := generateBatchID()

	req := &pb.BatchPDFProcessRequest{
		BucketName: bu.bucketName,
		BatchId:    batchId,
		UserId:     userId,
		PdfFiles:   UploadedFiles,
	}

	grpcRes, err := bu.client.ProcessBatchPDF(ctx, req)
	if err != nil {
		return nil, err
	}

	res := bu.mappingGrpcResponse(grpcRes)

	return res, nil
}

func (bu *BatchPdfProcessingUsecase) FetchSummarizedPdfHistory(ctx context.Context, userId string) (*model.CVSummarizeResponse, error) {
	req := &pb.FetchSummarizedPdfHistoryRequest{
		UserId: userId,
	}

	grpcRes, err := bu.client.FetchSummarizedPdfHistory(ctx, req)
	if err != nil {
		return nil, err
	}

	res := bu.mappingGrpcResponse(grpcRes)

	return res, nil
}

func (bu *BatchPdfProcessingUsecase) processPDF(ctx context.Context, zipFile *zip.File) (*pb.PDFFileInfo, error) {
	fileInZip, err := zipFile.Open()
	if err != nil {
		return nil, err
	}
	defer fileInZip.Close()

	fileContent, err := io.ReadAll(fileInZip)
	if err != nil {
		return nil, err
	}

	fileName := bu.generateUniqueFileName(zipFile.Name)

	reader := bytes.NewReader(fileContent)

	if err = bu.minio.MakeBucket(ctx, bu.bucketName); err != nil {
		return nil, err
	}

	info, err := bu.minio.UploadPDF(ctx,
		bu.bucketName, fileName, reader, int64(len(fileContent)))
	if err != nil {
		return nil, err
	}

	return &pb.PDFFileInfo{
		FileName:   fileName,
		Size:       info.Size,
		UploadedAt: timestamppb.New(time.Now()),
	}, nil
}

func (bu *BatchPdfProcessingUsecase) mappingGrpcResponse(res *pb.BatchPDFProcessResponse) *model.CVSummarizeResponse {
	response := &model.CVSummarizeResponse{}

	for _, pred := range res.Predictions {
		if pred.Prediction == nil {
			continue
		}
		response.Name = append(response.Name, pred.Prediction.Name...)
		response.CollegeName = append(response.CollegeName, pred.Prediction.CollegeName...)
		response.Degree = append(response.Degree, pred.Prediction.Degree...)
		response.GraduationYear = append(response.GraduationYear, pred.Prediction.GraduationYear...)
		response.YearsOfExperience = append(response.YearsOfExperience, pred.Prediction.YearsOfExperience...)
		response.CompaniesWorkedAt = append(response.CompaniesWorkedAt, pred.Prediction.CompaniesWorkedAt...)
		response.Designation = append(response.Designation, pred.Prediction.Designation...)
		response.Skills = append(response.Skills, pred.Prediction.Skills...)
		response.Location = append(response.Location, pred.Prediction.Location...)
		response.EmailAddress = append(response.EmailAddress, pred.Prediction.EmailAddress...)
	}

	return response
}

func (bu *BatchPdfProcessingUsecase) generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)

	baseName := strings.TrimSuffix(filepath.Base(originalName), ext)

	baseName = sanitizeFileName(baseName)

	uniqueID := uuid.NewString()

	timestamp := time.Now().Format("20060102-150405")

	return fmt.Sprintf("%s-%s-%s%s", baseName, timestamp, uniqueID[:8], ext)
}

func sanitizeFileName(name string) string {
	name = strings.ReplaceAll(name, " ", "_")

	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, name)

	return name
}

func generateBatchID() string {
	timestamp := time.Now().Format("20060102-150405")
	randomID := uuid.NewString()[:8]
	return fmt.Sprintf("cv-batch-%s-%s", timestamp, randomID)
}
