package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

func WinnersRegisterHandlers(r *mux.Router) {
  r.HandleFunc("/winners", WinnersUpdateHandler).Methods("POST")
}

func WinnersUpdateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse token info
  claims, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  // Check admin
  if !claims.Admin {
    w.WriteHeader(http.StatusUnauthorized)
    SendJson(w, JsonError{ Error: "Must be an admin" })
    log.Print("Must be an admin")
    return
  }

  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body WinnersUpdateRequest
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
  _, err = db.Exec("UPDATE winners SET winners = $1 WHERE id = 1", string(jsonBodyBytes))
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error updating picks" })
    log.Printf("Error updating picks: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)

  // Update sockets
  message, _ := json.Marshal(&WinnersMessage{
    Type: "winners",
    Winners: jsonBodyBytes,
  })

  // Send winners over as a message
  winnersHub.broadcast <- message
}

type WinnersUpdateRequest map[string]string
