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
  // Start connection
  client := serveWs(winnersHub, w, r, winnersReadHandler)

  // Get winners
  var winners string
  err := db.QueryRow("SELECT winners FROM winners WHERE id = 1").Scan(&winners)
  if err != nil {
    sendError(client, "Error querying database")
    log.Printf("Error querying database: %v", err)
    return
  }

  // Send over initial message
  sendWinners(client, winners)
}

func winnersReadHandler(client *SocketClient, message []byte) {
  log.Print(string(message))
}

func sendWinners(client *SocketClient, winners string) {
  // Convert to json
  message, _ := json.Marshal(&WinnersMessage{
    Type: "winners",
    Winners: []byte(winners),
  })

  // Send winners over as a message
  client.send <- message
}

type WinnersMessage struct {
  Type string `json:"type"`
  Winners json.RawMessage `json:"winners"`
}

func sendError(client *SocketClient, errMessage string) {
  message, _ := json.Marshal(&ErrorMessage{
    Type: "winners",
    Error: errMessage,
  })

  // Send winners over as a message
  client.send <- message
}

type ErrorMessage struct {
  Type string `json:"type"`
  Error string `json:"error"`
}

func SocketsChatHandler(w http.ResponseWriter, r *http.Request) {
  // Get messages
  rows, err := db.Query("SELECT messages.id, messages.user_id, users.email, messages.body, messages.created_at FROM messages INNER JOIN users ON users.id = messages.user_id")
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  var messages = []*ChatMessage{}

  // Iterate over messages
  for rows.Next() {
    var message ChatMessage

    err := rows.Scan(&message.Id, &message.UserId, &message.UserEmail, &message.Body, &message.CreatedAt)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      SendJson(w, JsonError{ Error: "Error querying database" })
      log.Printf("Error querying database: %v", err)
      return
    }

    messages = append(messages, &message)
  }

  // Check iteration for errors
  if err := rows.Err(); err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Convert to json
  response := &NewChatMessagesMessage{
    Type: "newChatMessages",
    ChatMessages: messages,
  }
  messagesBytes, err := json.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    SendJson(w, JsonError{ Error: "Error querying database" })
    log.Printf("Error querying database: %v", err)
    return
  }

  // Start connection
  client := serveWs(chatHub, w, r, chatReadHandler)

  // Send messages
  client.send <- messagesBytes
}

type NewChatMessagesMessage struct {
  Type string `json:"type"`
  ChatMessages []*ChatMessage `json:"chatMessages"`
}

type ChatMessage struct {
  Id int64 `json:"id"`
  UserId int64 `json:"userId"`
  UserEmail string `json:"email"`
  Body string `json:"body"`
  CreatedAt string `json:"createdAt"`
}

func chatReadHandler(client *SocketClient, message []byte) {
  var jsonMessage map[string]interface{}
  err := json.Unmarshal(message, &jsonMessage)
  if err != nil {
    message, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Can't parse message",
    })
    client.send <- message
    log.Printf("Can't parse message: %v", err)
    return
  }

  switch messageType := jsonMessage["type"]; messageType {
  case "postChatMessage":
    postChatMessage(client, jsonMessage)
  default:
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Invalid message type",
    })
    client.send <- errorMessage
    log.Printf("Invalid message type: %v", jsonMessage)
  }
}

func postChatMessage(client *SocketClient, message map[string]interface{}) {
  var token string
  var body string
  var ok bool

  // Validate token is string
  if token, ok = message["token"].(string); !ok {
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Missing token",
    })
    client.send <- errorMessage
    log.Printf("Missing token: %v", message)
    return
  }

  // Validate body is string
  if body, ok = message["body"].(string); !ok {
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Missing body",
    })
    client.send <- errorMessage
    log.Printf("Missing body: %v", message)
    return
  }


  // Get token claims
  claims, err := getAuthTokenClaimsFromString(token)
  if err != nil {
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Invalid token",
    })
    client.send <- errorMessage
    log.Printf("Invalid token: %v", message)
    return
  }

  // Create message
  var id int64
  var createdAt string
  err = db.QueryRow("INSERT INTO messages (user_id, body) VALUES ($1, $2) RETURNING id, created_at", claims.Id, body).Scan(&id, &createdAt)
  if err != nil {
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Error saving to database",
    })
    client.send <- errorMessage
    log.Printf("Error saving to database: %v", message)
    return
  }

  // Create responses
  responseMessage := &ChatMessage{
    Id: id,
    UserId: claims.Id,
    UserEmail: claims.Email,
    Body: body,
    CreatedAt: createdAt,
  }
  response := &NewChatMessagesMessage{
    Type: "newChatMessages",
    ChatMessages: []*ChatMessage{responseMessage},
  }
  messagesBytes, err := json.Marshal(response)
  if err != nil {
    errorMessage, _ := json.Marshal(&ErrorMessage{
      Type: "error",
      Error: "Error generating response",
    })
    client.send <- errorMessage
    log.Printf("Error generating response: %v", message)
    return
  }

  chatHub.broadcast <- messagesBytes
}
