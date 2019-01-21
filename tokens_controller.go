package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
  "os"
)

var tokenSecret string

func TokensRegisterHandlers(r *mux.Router) {
  // Ensure secret
  tokenSecret = os.Getenv("TOKEN_SECRET")
  if tokenSecret == "" {
    log.Fatal("Missing TOKEN_SECRET")
  }

  // Register routes
  r.HandleFunc("/tokens", TokensCreateHandler).Methods("POST")
}

func TokensCreateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body TokensCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Validate request
  if body.Email == "" || body.Password == "" {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }

  // Get user
  var id int64
  var passwordHash string
  var admin bool
  err = db.QueryRow("SELECT id, password_hash, admin FROM users WHERE email = $1", body.Email).Scan(&id, &passwordHash, &admin)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error finding user" })
    log.Print("Error finding user")
    return
  }

  // Check password
  if !ComparePasswords(body.Password, passwordHash) {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Wrong password" })
    log.Print("Wrong password")
    return
  }

  // Create token
  claims := &TokenClaims{
    Id: id,
    Email: body.Email,
    Admin: admin,
  }
  token, err := MakeToken(claims)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error creating token" })
    log.Printf("Error creating token: %s", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusCreated)
  SendJson(w, TokensCreateResponse{ Token: token })
}

type TokensCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

type TokensCreateResponse struct {
  Token string `json:"token"`
}
