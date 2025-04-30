module github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service

go 1.24.2

replace github.com/RecruitEase-Capstone/recruitEase-BE/authentication_service => ./

require (
	github.com/golang-jwt/jwt/v5 v5.2.2
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
	github.com/rs/zerolog v1.34.0
	golang.org/x/crypto v0.37.0
	google.golang.org/grpc v1.72.0
	google.golang.org/protobuf v1.36.6
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)
