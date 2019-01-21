package main

import (
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
