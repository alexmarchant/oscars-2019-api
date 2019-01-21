package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "strconv"
  "log"
  "fmt"
)

func PicksRegisterHandlers(r *mux.Router) {
  r.HandleFunc("/users/{id}/picks", PicksReadHandler).Methods("GET")
  r.HandleFunc("/users/{id}/picks", PicksUpdateHandler).Methods("POST")
}

func PicksReadHandler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  id := vars["id"]

  // Get picks
  var picks string
  err := db.QueryRow("SELECT picks FROM users WHERE id = $1", id).Scan(&picks)
  if err != nil {
    if err.Error() == "sql: no rows in result set" {
      w.WriteHeader(http.StatusNotFound)
      SendJson(w, JsonError{ Error: fmt.Sprintf("No user found with id %s", id) })
      log.Printf("No user found with id %s", id)
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
  vars := mux.Vars(r)
  id := vars["id"]

  // Parse id to int
  id64, err := strconv.ParseInt(id, 10, 64)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Invalid id" })
    log.Print("Invalid id")
    return
  }

  // Authorize user
  if !AuthorizeUser(r, id64) {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Not authorized to edit this user's picks" })
    log.Print("Not authorized to edit this user's picks")
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

  // Validate data
  valid := true
  for _, pick := range body {
    if pick.Category == "" || pick.Pick == "" {
      valid = false
    }
  }
  if !valid {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Invalid request body" })
    log.Print("Invalid request body")
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
  _, err = db.Exec("UPDATE users SET picks = $1 WHERE id = $2", string(jsonBodyBytes), id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error updating database" })
    log.Printf("Error updating database: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)
}

type PicksUpdateRequest []PicksUpdatePick

type PicksUpdatePick struct {
  Category string `json:"category"`
  Pick string `json:"pick"`
}
