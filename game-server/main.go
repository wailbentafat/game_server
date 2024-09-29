package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
)

var Upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type USER struct {
	user_id int `json:"user_id"`
	health  int `json:"health"`
}

type Move struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	playerID int `json:"playerID"`
}

type Shoot struct {
	Playerid int `json:"playerhurt"`
	damage   int `json:"damage"`
}

var clients = make(map[int]*websocket.Conn)
var mu sync.Mutex 


func sra_damage(player *USER, damage int) ([]byte) {
	player.health -= damage
	var message []byte
	if player.health <= 0 {
		player.health = 0
		message = []byte(fmt.Sprintf("User %d is dead", player.user_id))
	} else {
		message = []byte(fmt.Sprintf("User %d took %d damage", player.user_id, damage))
	}
	log.Printf("User %d took %d damage", player.user_id, damage)
	return message
}


func sra_move(player *USER, move Move) ([]byte) {
	var message []byte
	log.Printf("User %d moved to %d, %d", player.user_id, move.X, move.Y)
	message = []byte(fmt.Sprintf("User %d moved to %d, %d", player.user_id, move.X, move.Y))
	return message
}


func broadcastMessage(message []byte) {
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error sending message:", err)
			conn.Close() 
		}
	}
}


func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var user USER
	if err := conn.ReadJSON(&user); err != nil {
		log.Println("Error reading user:", err)
		return
	}

	
	mu.Lock()
	clients[user.user_id] = conn
	mu.Unlock()

	
	initialMessage := []byte("Hello, World!")
	if err := conn.WriteMessage(websocket.TextMessage, initialMessage); err != nil {
		log.Println("Error writing message:", err)
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
			continue
		}

		var message struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(p, &message); err != nil {
			continue
		}

		switch message.Type {
		case "move":
			var move Move
			if err := json.Unmarshal(p, &move); err != nil {
				continue
			}
			player := &USER{user_id: move.playerID, health: 100}
			msg := sra_move(player, move)
			broadcastMessage(msg)

		case "shoot":
			var shoot Shoot
			if err := json.Unmarshal(p, &shoot); err != nil {
				continue
			}
			player := &USER{user_id: shoot.Playerid, health: 100}
			msg := sra_damage(player, shoot.damage)
			broadcastMessage(msg)

		case "join":
			mu.Lock()
			clients[user.user_id] = conn
			mu.Unlock()
		}
	}
}

func main() {
	http.HandleFunc("/ws", wshandler)
	log.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
