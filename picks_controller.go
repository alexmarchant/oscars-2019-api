package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
  "fmt"
)

func PicksRegisterHandlers(r *mux.Router) {
  r.HandleFunc("/users/current-user/picks", PicksReadHandler).Methods("GET")
  r.HandleFunc("/users/current-user/picks", PicksUpdateHandler).Methods("POST")
}

func PicksReadHandler(w http.ResponseWriter, r *http.Request) {
  // Parse token info
  claims, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    SendJson(w, JsonError{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  // Get picks
  var picks string
  err = db.QueryRow("SELECT picks FROM users WHERE id = $1", claims.Id).Scan(&picks)
  if err != nil {
    if err.Error() == "sql: no rows in result set" {
      w.WriteHeader(http.StatusNotFound)
      SendJson(w, JsonError{ Error: fmt.Sprintf("Can't find user") })
      log.Printf("Can't find user: %v", err)
    } else {
      w.WriteHeader(http.StatusInternalServerError)
      SendJson(w, JsonError{ Error: "Error querying database" })
      log.Printf("Error querying database %v", err)
    }
    return
  }

  // Convert picks to bytes for json.RawMessage
  var picksJson json.RawMessage = []byte(picks)

  // Respond
  w.WriteHeader(http.StatusOK)
  SendJson(w, picksJson)
}

func PicksUpdateHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusUnauthorized)
  SendJson(w, JsonError{ Error: "Voting has been disabled" })
  log.Print("Voting has been disabled")
  return
}

type PicksUpdateRequest map[string]string
