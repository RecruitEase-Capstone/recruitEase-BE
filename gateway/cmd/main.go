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
	"github.com/RecruitEase-Capstone/recruitEase-BE/gateway/internal/usecase"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/joho/godotenv"
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

	AuthPort := os.Getenv("AUTH_SERVICE_PORT")
	AuthHost := os.Getenv("AUTH_SERVICE_HOST")
	AuthUrl := fmt.Sprintf("%s:%s", AuthHost, AuthPort)

	conn, err := grpc.NewClient(AuthUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal().Str("service", service).
			AnErr("failed to create a grpc client", err)
	}
	defer conn.Close()

	client := pb.NewAuthenticationServiceClient(conn)

	authUsecase := usecase.NewAuthClient(client)
	authHandler := handler.NewAuthHandler(authUsecase)

	r := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", GatewayPort),
		Handler: r,
	}

	handler.AuthRoutes(r, authHandler)

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
