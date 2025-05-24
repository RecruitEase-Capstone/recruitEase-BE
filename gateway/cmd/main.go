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

	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/config"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/handler"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/middleware"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils"
	minioUtils "github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/utils/minio"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/joho/godotenv"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	service      = "API gateway"
	shutdownTime = 30 * time.Second
)

func main() {
	godotenv.Load()

	config := config.Load()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log := zerolog.New(os.Stderr)

	authUrl := fmt.Sprintf("%s:%s", config.AuthHost, config.AuthPort)

	mlUrl := fmt.Sprintf("%s:%s", config.MLHost, config.MLPort)

	minioEndpoint := fmt.Sprintf("%s:%s", config.MinioHost, config.MinioPort)

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
		Creds:  credentials.NewStaticV4(config.MinioAccessKeyID, config.MinioSecretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatal().Str("service", service).AnErr("failed to create a minio client", err)
	}

	jwtPkg, err := jwt.NewJwt(config.JwtSecretKey, config.JwtExpired)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).Msg("cannot init jwt")
	}

	minio := minioUtils.NewMinio(minioClient)

	authUsecase := usecase.NewAuthUsecase(authClient)
	authHandler := handler.NewAuthHandler(authUsecase)

	batchProcessorUsecase := usecase.NewBatchPdfProcessing(minio, batchProcessorClient, config.MinioBucketName)
	batchProcessorHandler := handler.NewBatchPdfProcessingHandler(batchProcessorUsecase)

	middleware := middleware.NewMiddleware(jwtPkg, log)

	mux := http.NewServeMux()

	handler.AuthRoutes(mux, authHandler)
	handler.BatchProcessingRoutes(mux, batchProcessorHandler, middleware)

	mux.HandleFunc("/health-check", utils.HealthCheckHandler())

	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "https://recruitease-capstone.vercel.app"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With", "Origin"},
		AllowCredentials: true,
		MaxAge:           86400,
	}

	wrappedHandler := cors.New(corsOptions).Handler(
		middleware.LoggingMiddleware(middleware.RateLimiter(mux)))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", config.GatewayPort),
		Handler:           wrappedHandler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	gs := gracefullyShutdown(server, log)

	go func() {
		log.Info().Msgf("%s running on %s", service, config.GatewayPort)
		if err := server.ListenAndServe(); err != nil {
			log.Debug().AnErr("server error", err)
			log.Fatal().Str("service", service).AnErr("failed to serve the server", err)
		}
	}()

	gs()

	log.Info().Msgf("%s exited gracfully", service)
}

func gracefullyShutdown(s *http.Server, log zerolog.Logger) func() {
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
