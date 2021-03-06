package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

func UsersRegisterHandlers(r *mux.Router) {
  r.HandleFunc("/users", UsersCreateHandler).Methods("POST")
  r.HandleFunc("/users", UsersIndexHandler).Methods("GET")
  r.HandleFunc("/users/current-user", UsersCurrentUserHandler).Methods("GET")
  r.HandleFunc("/users/current-user", UsersPatchCurrentUserHandler).Methods("PATCH")
}

func UsersCreateHandler(w http.ResponseWriter, r *http.Request) {
  // Parse request
  decoder := json.NewDecoder(r.Body)
  var body UsersCreateRequest
  err := decoder.Decode(&body)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error parsing request" })
    log.Print("Error parsing request")
    return
  }

  // Validate request
  if body.Email == "" || body.Password == "" || body.PasswordConfirmation == "" {
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

  // Check password length
  if len(body.Password) < 6 {
    w.WriteHeader(http.StatusBadRequest)
    SendJson(w, JsonError{ Error: "Password must be at least 6 characters long" })
    log.Print("Password must be at least 6 characters long")
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

  // Hash pw
  passwordHash := HashAndSalt(body.Password)

  // Create user
  err = db.QueryRow("INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id", body.Email, passwordHash).Scan(&id)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error writing user data to database" })
    log.Printf("Error writing user data to database: %v", err)
    return
  }

  // Create JWT Token
  claims := &TokenClaims{
    Id: id,
    Email: body.Email,
    Admin: false,
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
  SendJson(w, UsersCreateResponse{ Token: token })
}

type UsersCreateRequest struct {
  Email string `json:"email"`
  Password string `json:"password"`
  PasswordConfirmation string `json:"passwordConfirmation"`
}

type UsersCreateResponse struct {
  Token string `json:"token"`
}

func UsersIndexHandler(w http.ResponseWriter, r *http.Request) {
  // Query users
  rows, err := db.Query("SELECT id, email, paid, picks FROM users")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  var users = []UsersIndexUser{}

  // Iterate over users
  for rows.Next() {
    var user UsersIndexUser
    var picks string

    err := rows.Scan(&user.Id, &user.Email, &user.Paid, &picks)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      SendJson(w, JsonError{ Error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
      return
    }

    // Convert picks to bytes for json.RawMessage
    user.Picks = []byte(picks)
    users = append(users, user)
  }

  // Check iteration for errors
  if err := rows.Err(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)
  SendJson(w, users)
}

type UsersIndexUser struct {
  Id int64 `json:"id"`
  Email string `json:"email"`
  Paid bool `json:"paid"`
  Picks json.RawMessage `json:"picks"`
}

type UsersIndexResponse []UsersIndexUser

func UsersCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
  // Parse token info
  claims, err := getAuthTokenClaims(r)
  if err != nil {
    w.WriteHeader(http.StatusUnauthorized)
    SendJson(w, JsonError{ Error: "Invalid token" })
    log.Printf("Invalid token: %v", err)
    return
  }

  // Get user
  var user UsersCurrentUserResponse
  err = db.QueryRow("SELECT id, email, admin, paid FROM users WHERE id = $1", claims.Id).Scan(&user.Id, &user.Email, &user.Admin, &user.Paid)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Can't find user" })
    log.Printf("Can't find user: %v", err)
    return
  }

  // Respond
  w.WriteHeader(http.StatusOK)
  SendJson(w, user)
}

type UsersCurrentUserResponse struct {
  Id int64 `json:"id"`
  Email string `json:"email"`
  Admin bool `json:"admin"`
  Paid bool `json:"paid"`
}

func UsersPatchCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
  w.WriteHeader(http.StatusUnauthorized)
  SendJson(w, JsonError{ Error: "Voting has been disabled" })
  log.Print("Voting has been disabled")
  return
}

type PatchCurrentUserRequest struct {
  Paid bool `json:"paid"`
}
