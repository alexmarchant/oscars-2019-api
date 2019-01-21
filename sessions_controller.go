package main

import (
  "github.com/gorilla/mux"
  "net/http"
)

func RegisterSessionHandlers(r *mux.Router) {
  r.HandleFunc("/sessions", CreateSessionHandler).Methods("POST")
}

func CreateSessionHandler(w http.ResponseWriter, r *http.Request) {
}
