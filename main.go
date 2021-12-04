package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"

	"./server"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":443", config)
	if err != nil {
		log.Println(err)
		return
	}

	defer ln.Close()
	router := mux.NewRouter().StrictSlash(true)
	//ConsUSer.ConsultaUsuarios()
	server.AllRooms.Init()
	// do all your routes declaration

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})
	router.HandleFunc("/hello", server.HelloWord)

	router.HandleFunc("/create", server.CreateRoomRequestHandler)
	router.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.ListenAndServeTLS(":"+port, "server.crt", "server.key", handlers.CORS(headers, origins, methods)(router))

}
