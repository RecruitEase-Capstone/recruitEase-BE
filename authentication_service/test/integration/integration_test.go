package integration_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"buf.build/go/protovalidate"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/handler"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/repository"
	"github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/usecase"
	"github.com/RecruitEase-Capstone/recruitEase-BE/pkg/jwt"
	pb "github.com/RecruitEase-Capstone/recruitEase-BE/pkg/proto/v1"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	customErr "github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service/internal/utils/error"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	tcPostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName               = "db-test"
	dbUser               = "test"
	dbPassword           = "passtest"
	postgresInternalPort = "5432/tcp"
	secretKey            = "6af77e29d3a2d2d2215da854e1a81257a1b844a018eff71a0486f8f4395ebf7e"
	jwtExpired           = "30m"
)

type IntegrationTestSuite struct {
	suite.Suite
	postgresContainer testcontainers.Container
	DB                *sqlx.DB
	AuthRepo          repository.IAuthRepository
	AuthUsecase       usecase.IAuthUsecase
	AuthHandler       *handler.AuthHandler
	ctx               context.Context
}

func (i *IntegrationTestSuite) SetupSuite() {
	i.ctx = context.Background()

	postgresContainer, err := i.setupPostgresContainer()
	if err != nil {
		i.T().Fatalf("failed to setup postgres container: %v", err)
	}
	i.postgresContainer = postgresContainer

	db, err := i.connectToPostgres()
	if err != nil {
		i.T().Fatalf("failed connect to postgres: %v", err)
	}
	i.DB = db

	err = i.runMigration()
	if err != nil {
		i.T().Fatalf("failed to run migration: %v", err)
	}

	err = i.setupLayers()
	if err != nil {
		i.T().Fatalf("failed to setup layers: %v", err)
	}
}

func (i *IntegrationTestSuite) SetupTest() {
	if i.DB != nil {
		_, err := i.DB.Exec("TRUNCATE TABLE users")
		assert.NoError(i.T(), err)
	}
}

func (i *IntegrationTestSuite) setupPostgresContainer() (testcontainers.Container, error) {
	ctx := i.ctx

	postgresContainer, err := tcPostgres.Run(ctx,
		"postgres:17",
		tcPostgres.WithDatabase(dbName),
		tcPostgres.WithUsername(dbUser),
		tcPostgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(120*time.Second),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create postgres container: %w", err)
	}

	return postgresContainer, nil
}

func (i *IntegrationTestSuite) TearDownSuite() {
	if i.DB != nil {
		if err := i.DB.Close(); err != nil {
			log.Info().Msgf("failed to close database connection: %v", err)
		}
	}

	if i.postgresContainer != nil {
		if err := i.postgresContainer.Terminate(i.ctx); err != nil {
			log.Info().Msgf("failed to terminate container: %v", err)
		}
	}
}

func (i *IntegrationTestSuite) connectToPostgres() (*sqlx.DB, error) {
	ctx := i.ctx
	container := i.postgresContainer

	if container == nil {
		return nil, fmt.Errorf("postgres container is nil")
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(postgresInternalPort))
	if err != nil {
		return nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get container host: %w", err)
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, host, mappedPort.Port(), dbName)

	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres : %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return db, nil
}

func (i *IntegrationTestSuite) setupLayers() error {
	i.AuthRepo = repository.NewAuthRepository(i.DB)

	jwt, err := jwt.NewJwt(secretKey, jwtExpired)
	if err != nil {
		return fmt.Errorf("failed to create JWT: %w", err)
	}

	i.AuthUsecase = usecase.NewAuthUsecase(i.AuthRepo, jwt)

	validator, err := protovalidate.New()
	if err != nil {
		return fmt.Errorf("failed to create validator: %w", err)
	}

	i.AuthHandler = handler.NewAuthHandler(i.AuthUsecase, validator)

	return nil
}

func (i *IntegrationTestSuite) runMigration() error {
	if i.DB == nil {
		return fmt.Errorf("database connection is nil")
	}

	migrationFilePath := "file://../../db/migrations"

	driver, err := postgres.WithInstance(i.DB.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create postgres driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationFilePath,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

// Test cases
func (i *IntegrationTestSuite) TestRegisterUser() {
	testCases := []struct {
		name          string
		setup         func() error
		input         *pb.RegisterRequest
		expectedName  string
		expectedEmail string
		expectedError error
	}{
		{
			name: "Success - Valid Registration",
			input: &pb.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "Rahasia#123",
				ConfirmPassword: "Rahasia#123",
			},
			expectedName:  "John Doe",
			expectedEmail: "john@example.com",
			expectedError: nil,
		},
		{
			name: "Failed - Password too short",
			input: &pb.RegisterRequest{
				Name:            "John Doe",
				Email:           "john@example.com",
				Password:        "Pendek",
				ConfirmPassword: "Pendek",
			},
			expectedError: customErr.ErrPasswordTooShort,
		},
		{
			name: "Failed - Email is exists",
			setup: func() error {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Rahasia#123"), bcrypt.DefaultCost)
				assert.NoError(i.T(), err)

				_, err = i.DB.ExecContext(i.ctx,
					"INSERT INTO users (id, name, email, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)",
					uuid.NewString(),
					"First John Doe",
					"emailExists@example.com",
					string(hashedPassword),
					time.Now(),
					time.Now(),
				)
				assert.NoError(i.T(), err)

				return nil
			},
			input: &pb.RegisterRequest{
				Name:            "Second John Doe",
				Email:           "emailExists@example.com",
				Password:        "Rahasia#123",
				ConfirmPassword: "Rahasia#123",
			},
			expectedError: customErr.ErrEmailExist,
		},
	}

	for _, tc := range testCases {
		i.T().Run(tc.name, func(t *testing.T) {

			if tc.setup != nil {
				err := tc.setup()
				assert.NoError(t, err, "Setup failed")
			}

			res, err := i.AuthHandler.UserRegister(i.ctx, tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Equal(t, tc.expectedName, res.Name)
				assert.Equal(t, tc.expectedEmail, res.Email)
			}
		})
	}

}

func (i *IntegrationTestSuite) TestLoginUser() {
	req := &pb.RegisterRequest{
		Name:            "John Doe",
		Email:           "john@example.com",
		Password:        "Rahasia#123",
		ConfirmPassword: "Rahasia#123",
	}

	_, err := i.AuthHandler.UserRegister(i.ctx, req)
	assert.NoError(i.T(), err)

	testCases := []struct {
		name          string
		input         *pb.LoginRequest
		expectedError error
	}{
		{
			name: "Success - Valid Credentials",
			input: &pb.LoginRequest{
				Email:    "john@example.com",
				Password: "Rahasia#123",
			},
			expectedError: nil,
		},
		{
			name: "Failed - invalid credentials",
			input: &pb.LoginRequest{
				Email:    "wrongemail@example.com",
				Password: "wrongPw123",
			},
			expectedError: customErr.ErrInvalidCredentials,
		},
	}

	for _, tc := range testCases {
		i.T().Run(tc.name, func(t *testing.T) {
			res, err := i.AuthHandler.UserLogin(i.ctx, tc.input)

			if tc.expectedError != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.NotEmpty(t, res.JwtToken)
			}
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
