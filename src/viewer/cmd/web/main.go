package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"strings"
	"viewer/pkg/env"
	"viewer/pkg/page"
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
	mux.HandleFunc("/", handleIndex())
	mux.Handle("/public/", handlePublic("/public/"))
	// mux.HandleFunc("/socket", socketHandler)

	http.ListenAndServe(":"+*port, mux)
}

func handleIndex() http.HandlerFunc {
	param := page.StatusParams{
		IndexParams: page.IndexParams{StylesheetURL: "/public/style.css"},
		Title:       "Status viewer",
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("index handler called for", r.URL)
		buf := &bytes.Buffer{}
		err := page.Status(buf, param)
		if err != nil {
			log.Printf("failed to generate status page: %+v", err)
			http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)
	}
}

func handlePublic(path string) http.Handler {
	log.Println("serving files at", path)
	// Necessary to strip because file server serves relative to ./public/ folder.
	return http.StripPrefix(path,
		http.FileServer(http.Dir(strings.ReplaceAll(path, "/", ""))))
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
