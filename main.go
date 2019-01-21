package main

import (
  "io"
  "github.com/gorilla/mux"
  "net/http"
)

func main() {
  r := mux.NewRouter()

  // DB
  ConnectDB()

  // Routes
  r.HandleFunc("/", HomeHandler)
  RegisterUserHandlers(r)
  RegisterTokenHandlers(r)

  // Start server
  http.Handle("/", r)
  http.ListenAndServe(":3000", r)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusOK)
  io.WriteString(w, "OK")
}
