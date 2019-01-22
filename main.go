package main

import (
  "github.com/gorilla/handlers"
  "github.com/gorilla/mux"
  "net/http"
  "io"
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

  // Start server
  http.Handle("/", r)
  http.ListenAndServe(":3000", handlers.CORS()(r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "OK")
}
