package main

import (
  "github.com/dgrijalva/jwt-go"
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
  "os"
)

var tokenSecret string

func RegisterTokenHandlers(r *mux.Router) {
  // Ensure secret
  tokenSecret = os.Getenv("TOKEN_SECRET")
  if tokenSecret == "" {
    log.Fatal("Missing TOKEN_SECRET")
  }

  // Register routes
  r.HandleFunc("/tokens", CreateTokenHandler).Methods("POST")
}

func CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body TokenCreateRequest
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
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
      "email": body.Email,
      "admin": admin,
  })
  tokenString, err := token.SignedString([]byte(tokenSecret))

  // Return token
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error creating token" })
    log.Printf("Error creating token: %s", err)
    return
  }

  w.WriteHeader(http.StatusCreated)
  SendJson(w, TokenCreateResponse{ Token: tokenString })
}

type TokenCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

type TokenCreateResponse struct {
  Token string `json:"token"`
}
