package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//AllRooms is the global hashmap for the server
var AllRooms RoomMap

func HelloWord(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)
	json.NewEncoder(w).Encode("message")
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}
	log.Println(AllRooms)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})

}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Fatal(err)
					client.Conn.Close()

				}

			}
			fmt.Println(client.Conn)
		}
	}

}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]
	if !ok {
		log.Println("roomId missing in URL Parameters")
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("web socket upgrade Error", err)
	}
	AllRooms.InsertIntoRoom(roomID[0], false, ws)
	go broadcaster()

	for {
		var msg broadcastMsg

		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("read error", err)
		}
		msg.Client = ws
		msg.RoomID = roomID[0]

		broadcast <- msg
	}
}
