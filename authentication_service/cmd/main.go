package main

import (
	"fmt"
	"net"
	"os"

	"buf.build/go/protovalidate"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/db"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/handler"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/repository"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/usecase"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/jwt"
	log_utils "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/log"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	logging_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

const service = "authentication service"

func main() {
	godotenv.Load()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

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

	logger, opts := log_utils.InterceptorLogger(zerolog.New(os.Stderr))

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

	log.Info().Msgf("%s running on %s", service, AuthPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Err(err).Str("service", service).Msgf("failed to serve %s", service)
	}
}
