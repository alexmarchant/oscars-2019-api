package main

import (
  "github.com/rs/cors"
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
  SocketsRegisterHandlers(r)

  // Start server
  http.ListenAndServe(":3000", cors.Default().Handler(r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "OK")
}
