package main

import (
  "github.com/dgrijalva/jwt-go"
  "golang.org/x/crypto/bcrypt"
  "log"
)

func HashAndSalt(password string) string {
  passwordBytes := []byte(password)
  hash, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
  if err != nil {
    log.Println(err)
  }
  return string(hash)
}

func ComparePasswords(password string, passwordHash string) bool {
  passwordBytes := []byte(password)
  hashBytes := []byte(passwordHash)
  err := bcrypt.CompareHashAndPassword(hashBytes, passwordBytes)
  if err != nil {
      log.Println(err)
      return false
  }

  return true
}

func MakeToken(id int64, email string, admin bool) (string, error) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    "id": id,
    "email": email,
    "admin": admin,
  })
  return token.SignedString([]byte(tokenSecret))
}
