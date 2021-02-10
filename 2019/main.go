package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/mcraealex/BattleSnake2019/routes"
	"github.com/mcraealex/BattleSnake2019/server"
)

var (
	serverAddress = ":80"
	certFile      = "./cert.crt"
	keyFile       = "./server.key"
	databasename  = "./testing.db"
)

func main() {
	logger := log.New(os.Stderr, "Testing: ", log.LstdFlags|log.Lshortfile)

	// setup router
	h := routes.NewHandlers(logger) //, dbh
	gmux := mux.NewRouter()
	h.SetupRoutes(gmux)
	// setup server
	srv := server.New(gmux, serverAddress)
	// run server
	err := srv.ListenAndServe()
	if err != nil {
		logger.Fatalf("Error starting server: %v\n", err)
	}
}
