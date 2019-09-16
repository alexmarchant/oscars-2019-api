package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	r := mux.NewRouter()

	// DB
	ConnectDB()

	// Routes
	r.HandleFunc("/", HomeHandler)
	UsersRegisterHandlers(r)
	TokensRegisterHandlers(r)
	PicksRegisterHandlers(r)
	WinnersRegisterHandlers(r)
	SocketsRegisterHandlers(r)
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// CORS
	cors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{"*"},
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	portString := fmt.Sprintf(":%s", port)
	http.ListenAndServe(portString, cors.Handler(r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "OK")
}
