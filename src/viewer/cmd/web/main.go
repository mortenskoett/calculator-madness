package main

import (
	"flag"
	"log"
	"net/http"
	"viewer/pkg/env"
)

const (
	ServiceNameChannel string = "viewer"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
	port           = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
)

// var upgrader = websocket.Upgrader{}

func main() {

	log.Println("starting calculator viewer web service")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	// mux.HandleFunc("/socket", socketHandler)

	http.ListenAndServe(":"+*port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Hello Anne Dorte!</h1>"))
}

// func socketHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Print("attempting upgrade to websocket")
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Print("upgrade failed: ", err)
// 		return
// 	}
// 	defer conn.Close()

// 	mt, message, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("read failed:", err)
// 	}
// 	log.Print(string(message))

// 	retmsg := []byte("<h1>Hello Anne Dorte!</h1>")

// 	err = conn.WriteMessage(mt, retmsg)
// }
