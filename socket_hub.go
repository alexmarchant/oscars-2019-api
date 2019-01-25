// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type SocketHub struct {
  // Registered clients.
  clients map[*SocketClient]bool

  // Inbound messages from the clients.
  broadcast chan []byte

  // Register requests from the clients.
  register chan *SocketClient

  // Unregister requests from clients.
  unregister chan *SocketClient
}

func newHub() *SocketHub {
  return &SocketHub{
    broadcast:  make(chan []byte),
    register:   make(chan *SocketClient),
    unregister: make(chan *SocketClient),
    clients:    make(map[*SocketClient]bool),
  }
}

func (h *SocketHub) run() {
  for {
    select {
    case client := <-h.register:
      h.clients[client] = true
    case client := <-h.unregister:
      if _, ok := h.clients[client]; ok {
        delete(h.clients, client)
        close(client.send)
      }
    case message := <-h.broadcast:
      for client := range h.clients {
        select {
        case client.send <- message:
        default:
          close(client.send)
          delete(h.clients, client)
        }
      }
    }
  }
}
