package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func DBConnection() *sqlx.DB {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, username, password, database, port)

	log.Print("connecting to database...")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal().
			Err(err).Msg("failed to connect to postgres...")
	}
	log.Print("success connecting to database...")

	return db
}
