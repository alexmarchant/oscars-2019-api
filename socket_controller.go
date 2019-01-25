package main

import (
  "github.com/gorilla/mux"
  "net/http"
)

var winnersHub *SocketHub

func SocketsRegisterHandlers(r *mux.Router) {
  winnersHub = newHub()
  go winnersHub.run()
  r.HandleFunc("/ws/winners", SocketsWinnersHandler)
}

func SocketsWinnersHandler(w http.ResponseWriter, r *http.Request) {
  serveWs(winnersHub, w, r)
}
