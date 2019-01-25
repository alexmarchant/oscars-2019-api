package main

import (
  "github.com/gorilla/mux"
  "encoding/json"
  "net/http"
  "log"
)

var winnersHub *SocketHub
var chatHub *SocketHub

func SocketsRegisterHandlers(r *mux.Router) {
  // Winners
  winnersHub = newHub()
  go winnersHub.run()
  r.HandleFunc("/ws/winners", SocketsWinnersHandler)

  // Chat
  chatHub = newHub()
  go chatHub.run()
  r.HandleFunc("/ws/chat", SocketsChatHandler)
}

func SocketsWinnersHandler(w http.ResponseWriter, r *http.Request) {
  // Get winners
  var winners string
  err := db.QueryRow("SELECT winners FROM winners WHERE id = 1").Scan(&winners)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Start connection
  serveWs(winnersHub, w, r, winnersReadHandler)

  // Send winners
  winnersHub.broadcast <- []byte(winners)
}

func winnersReadHandler(message []byte) {
  log.Print(string(message))
}

func SocketsChatHandler(w http.ResponseWriter, r *http.Request) {
  // Get messages
  rows, err := db.Query("SELECT id, user_id, body, created_at FROM messages")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  var messages = []Message{}

  // Iterate over messages
  for rows.Next() {
    var message Message

    err := rows.Scan(&message.Id, &message.UserId, &message.CreatedAt, &message.Body)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      SendJson(w, JsonError{ Error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
      return
    }

    messages = append(messages, message)
  }

  // Check iteration for errors
  if err := rows.Err(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Convert to json
  messagesBytes, err := json.Marshal(messages)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Start connection
  serveWs(chatHub, w, r, chatReadHandler)

  // Send messages
  chatHub.broadcast <- messagesBytes
}

type Message struct {
  Id int64 `json:"id"`
  UserId int64 `json:"userId"`
  Body string `json:"body"`
  CreatedAt string `json:"createdAt"`
}

func chatReadHandler(message []byte) {
  log.Print(string(message))
}
