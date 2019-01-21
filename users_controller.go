package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

func RegisterUserHandlers(r *mux.Router) {
  r.HandleFunc("/users", CreateUserHandler).Methods("POST")
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body UserCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Missing required params" })
    log.Print("Missing required params")
    return
  }

  // Check password matched passwordConfirmation
  if body.Password != body.PasswordConfirmation {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Password doesn't match confirmation" })
    log.Print("Password doesn't match confirmation")
    return
  }

  // Check if user exists
  var id int64
  err = db.QueryRow("SELECT id FROM users WHERE email = $1", body.Email).Scan(&id)
  if err == nil {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "User already exists" })
    log.Print("User already exists")
    return
  }

  passwordHash := HashAndSalt(body.Password)

  // Create user
  err = db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", body.Email, passwordHash).Scan(&id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error writing user data to database" })
    log.Printf("Error writing user data to database: %v", err)
    return
  }

  w.WriteHeader(http.StatusCreated)
  SendJson(w, UserCreatedResponses{ Id: id })
}

type UserCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
  PasswordConfirmation string `json:"passwordConfirmation"`
}

type UserCreatedResponses struct {
  Id int64 `json:"id"`
}
