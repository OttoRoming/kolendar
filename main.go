package main

import (
	"log"
	"os"

	_ "embed"

	"github.com/OttoRoming/kolendar/server"
)

var (
	//go:embed schema.sql
	schema string
)

func main() {
	server, err := server.NewServer(schema)
	if err != nil {
		log.Fatal("failed to start server", "error", err)
		os.Exit(1)
	}

	if err := server.Run(); err != nil {
		log.Fatal("server error", "error", err)
		os.Exit(1)
	}

	server.Close()
}
