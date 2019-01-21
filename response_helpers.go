package main

import (
  "net/http"
  "encoding/json"
)

type JsonError struct {
  Error string `json:"error"`
}

func SendJson(w http.ResponseWriter, data interface{}) {
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(data)
}
