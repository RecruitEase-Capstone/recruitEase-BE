package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/handler"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/middleware"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	minioUtils "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/minio"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	service      = "API gateway"
	shutdownTime = 30 * time.Second
)

func main() {
	godotenv.Load()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	GatewayPort := os.Getenv("GATEWAY_SERVICE_PORT")

	authPort := os.Getenv("AUTH_SERVICE_PORT")
	authHost := os.Getenv("AUTH_SERVICE_HOST")
	authUrl := fmt.Sprintf("%s:%s", authHost, authPort)

	mlPort := os.Getenv("ML_SERVICE_PORT")
	mlHost := os.Getenv("ML_SERVICE_HOST")
	mlUrl := fmt.Sprintf("%s:%s", mlHost, mlPort)

	minioHost := os.Getenv("MINIO_API_HOST")
	minioPort := os.Getenv("MINIO_API_PORT")
	minioEndpoint := fmt.Sprintf("%s:%s", minioHost, minioPort)
	log.Info().Str("minio url", minioEndpoint)

	minioAccessKeyID := os.Getenv("MINIO_ROOT_USER")
	minioSecretAccessKey := os.Getenv("MINIO_ROOT_PASSWORD")
	JwtKey := os.Getenv("JWT_SECRET_KEY")
	JwtExpired := os.Getenv("JWT_EXPIRED")
	bucketName := os.Getenv("MINIO_BUCKET_NAME")

	authConn, err := grpc.NewClient(authUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Str("service", service).
			AnErr("failed to create a authentication service grpc client", err)
	}
	defer authConn.Close()

	mlConn, err := grpc.NewClient(mlUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Str("service", service).
			AnErr("failed to create a ml service grpc client", err)
	}
	defer mlConn.Close()

	authClient := pb.NewAuthenticationServiceClient(authConn)
	batchProcessorClient := pb.NewCVProcessorServiceClient(mlConn)

	minioClient, err := minio.New(minioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKeyID, minioSecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Str("service", service).AnErr("failed to create a minio client", err)
	}

	jwtPkg, err := jwt.NewJwt(JwtKey, JwtExpired)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).Msg("cannot init jwt")
	}

	minio := minioUtils.NewMinio(minioClient)

	authUsecase := usecase.NewAuthClient(authClient)
	authHandler := handler.NewAuthHandler(authUsecase)

	batchProcessorUsecase := usecase.NewBatchPdfProcessing(minio, batchProcessorClient, bucketName)
	batchProcessorHandler := handler.NewBatchPdfProcessingHandler(batchProcessorUsecase)

	middleware := middleware.NewMiddleware(jwtPkg)

	mux := http.NewServeMux()

	handler.AuthRoutes(mux, authHandler)
	handler.BatchProcessingRoutes(mux, batchProcessorHandler, middleware)

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "Origin"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	corsMiddleware := cors.New(corsOptions).Handler(mux)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", GatewayPort),
		Handler: corsMiddleware,
	}

	gs := gracefullyShutdown(server)

	go func() {
		log.Info().Msgf("%s running on %s", service, GatewayPort)
		if err := server.ListenAndServe(); err != nil {
			log.Debug().AnErr("server error", err)
			log.Fatal().Str("service", service).AnErr("failed to serve the server", err)
		}
	}()

	gs()

	log.Info().Msgf("%s exited gracfully", service)
}

func gracefullyShutdown(s *http.Server) func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return func() {
		sig := <-c

		log.Info().Msgf("received shutdown signal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTime)
		defer cancel()

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()

			log.Info().Msgf("Gracefully stopping %s... ", service)

			stopped := make(chan struct{})
			go func() {
				s.Shutdown(ctx)
				close(stopped)
			}()

			select {
			case <-stopped:
				log.Info().Msg("All connections completed gracefully")
			case <-ctx.Done():
				log.Info().Msg("Shutdown time reached, forcing stop")
				s.Close()
			}
		}()

		wg.Wait()
	}
}
