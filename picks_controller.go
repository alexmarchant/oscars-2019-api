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
  // Parse token info
  claims, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    SendJson(w, JsonError{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body PicksUpdateRequest
  err = decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Turn request into string for postgres
  jsonBodyBytes, err := json.Marshal(body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Update postgres data
  _, err = db.Exec("UPDATE users SET picks = $1 WHERE id = $2", string(jsonBodyBytes), claims.Id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error updating picks" })
    log.Printf("Error updating picks: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)
}

type PicksUpdateRequest map[string]string
