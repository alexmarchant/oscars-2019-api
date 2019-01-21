package main

import (
  "github.com/gorilla/mux"
  "net/http"
  "log"
)

func WinnersRegisterHandlers(r *mux.Router) {
  r.HandleFunc("/winners", WinnersUpdateHandler).Methods("POST")
}

func WinnersUpdateHandler(w http.ResponseWriter, r *http.Request) {
  if !AuthorizeAdmin(r) {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Must be an admin" })
    log.Print("Must be an admin")
    return
  }
}
