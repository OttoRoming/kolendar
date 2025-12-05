package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"

	"database/sql"

	_ "embed"

	"github.com/OttoRoming/kolendar/repository"
	"github.com/joho/godotenv"
	_ "modernc.org/sqlite"
)

var (
	//go:embed schema.sql
	ddl       string
	queries   *repository.Queries
	jwtSecret string
)

func setup() error {
	ctx := context.Background()

	err := godotenv.Load()
	if err != nil {
		return err
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return errors.New("JWT_SECRET not set in environment")
	}

	// create file database
	db, err := sql.Open("sqlite", "file:kolendar.db?_foreign_keys=on")
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries = repository.New(db)

	return nil
}

func run() error {
	if err := setup(); err != nil {
		return err
	}

	// run the http server
	router := http.NewServeMux()

	router.HandleFunc("POST /users/", createUser)
	router.HandleFunc("POST /users/login", loginUser)

	// get the users calendars
	router.HandleFunc("GET /calendars/", getCalendars)
	router.HandleFunc("POST /calendars/", createCalendar)

	err := http.ListenAndServe(":8000", router)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
