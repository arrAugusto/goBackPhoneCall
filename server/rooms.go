package server

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//Participant describes a single entityt in the hashmap
type Participant struct {
	Host bool
	Conn *websocket.Conn
}
type RoomMap struct {
	Mutext sync.RWMutex
	Map    map[string][]Participant
}

//init inicializando el roommpa struct

func (r *RoomMap) Init() {
	r.Map = make(map[string][]Participant)
}

//Get will return the array of participants int the rooms
func (r *RoomMap) Get(roomID string) []Participant {
	r.Mutext.RUnlock()
	defer r.Mutext.Lock()
	return r.Map[roomID]
}

//creando room
func (r *RoomMap) CreateRoom() string {
	r.Mutext.Lock()
	defer r.Mutext.Unlock()

	rand.Seed(time.Now().UnixNano())
	var letters = []rune("459567d3bde4418b7fe302ff9809c4b0befaf7dd")
	b := make([]rune, 8)

	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	roomID := string(b)
	r.Map[roomID] = []Participant{}
	return roomID
}

func (r *RoomMap) InsertIntoRoom(roomID string, host bool, conn *websocket.Conn) {
	r.Mutext.Lock()
	defer r.Mutext.Unlock()

	p := Participant{host, conn}

	log.Println("Insertando en el room con roomID", roomID)
	r.Map[roomID] = append(r.Map[roomID], p)

}

func (r *RoomMap) DeleteRoom(roomID string) {
	r.Mutext.Lock()
	defer r.Mutext.Unlock()
	delete(r.Map, roomID)
}
