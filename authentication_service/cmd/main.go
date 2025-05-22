package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"buf.build/go/protovalidate"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/db"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/handler"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/repository"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/usecase"
	log_utils "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/log"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
)

const (
	service      = "authentication service"
	shutdownTime = 30 * time.Second
)

func main() {
	godotenv.Load()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	log := zerolog.New(os.Stderr)

	AuthPort := fmt.Sprintf(":%s", os.Getenv("AUTH_SERVICE_PORT"))
	JwtKey := os.Getenv("JWT_SECRET_KEY")
	JwtExpired := os.Getenv("JWT_EXPIRED")

	db := db.DBConnection()
	defer db.Close()

	lis, err := net.Listen("tcp", AuthPort)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).Msg("failed to listen")
	}

	jwtPkg, err := jwt.NewJwt(JwtKey, JwtExpired)
	if err != nil {
		log.Fatal().
			Err(err).
			Str("service", service).Msg("cannot init jwt")
	}

	validator, err := protovalidate.New()
	if err != nil {
		log.Fatal().Err(err).Str("service", service).Msg("failed to initialize protovalidate")
	}

	logger, opts := log_utils.InterceptorLogger(log)

	protovalidateInterceptor := protovalidate_middleware.UnaryServerInterceptor(validator)

	loggingInterceptor := logging_middleware.UnaryServerInterceptor(logger, opts...)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(protovalidateInterceptor),
		grpc.ChainUnaryInterceptor(loggingInterceptor),
	)

	authRepo := repository.NewAuthRepository(db)
	authUsecase := usecase.NewAuthUsecase(authRepo, jwtPkg)
	authHanlder := handler.NewAuthHandler(authUsecase, validator)

	pb.RegisterAuthenticationServiceServer(grpcServer, authHanlder)

	gs := gracefullyShutdown(grpcServer, log)

	go func() {
		log.Info().Msgf("%s running on %s", service, AuthPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal().Err(err).Str("service", service).Msgf("failed to serve %s", service)
		}
	}()

	gs()

	log.Info().Msgf("%s exited gracefully", service)
}

func gracefullyShutdown(s *grpc.Server, log zerolog.Logger) func() {
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
				s.GracefulStop()
				close(stopped)
			}()

			select {
			case <-stopped:
				log.Info().Msg("All connections completed gracefully")
			case <-ctx.Done():
				log.Info().Msg("Shutdown time reached, forcing stop")
				s.Stop()
			}
		}()

		wg.Wait()
	}
}
