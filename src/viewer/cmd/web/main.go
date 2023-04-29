package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
	"viewer/pkg/env"
	"viewer/pkg/page"
)

const (
	ServiceNameChannel string = "viewer"
	indexStyleURL      string = "/public/indexstyle.css"
	statusStyleURL     string = "/public/statusstyle.css"
	faviconURL         string = "/public/images/calculator.ico"
)

var (
	nsqlookupdAddr = flag.String("nsqlookupd-addr", env.GetEnvVarOrDefault("NSQLOOKUPD_ADDR", "127.0.0.1:4161"), "Address of nsqlookupd server with port")
	port           = flag.String("port", env.GetEnvVarOrDefault("PORT", "3000"), "Port on which the UI will be served")
)

var (
	// FIXME: Dummy implementation. Delete later.
	calcs = []page.Calculation{
		{
			CalculationID: "1",
			MessageID:     "abc1",
			CreatedTime:   "1.1.1",
			Equation:      "1+1",
			Progress:      page.Progress{Current: 2, Outof: 5},
			Result:        "",
		},
		{
			CalculationID: "2",
			MessageID:     "abc2",
			CreatedTime:   "2.2.2",
			Equation:      "2+2",
			Progress:      page.Progress{Current: 0, Outof: 5},
			Result:        "",
		},
		{
			CalculationID: "3",
			MessageID:     "abc3",
			CreatedTime:   "3.3.3",
			Equation:      "3*3",
			Progress:      page.Progress{Current: 5, Outof: 5},
			Result:        "9",
		},
	}
)

// var upgrader = websocket.Upgrader{}

func main() {

	log.Println("starting calculator viewer web service")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", loggingHandlerFunc(handleIndex()))
	mux.HandleFunc("/favicon.ico", loggingHandlerFunc(serveFile(faviconURL)))
	mux.Handle("/public/", loggingHandler(fileServerHandler("/public/")))
	// mux.HandleFunc("/socket", socketHandler)
	http.ListenAndServe(":"+*port, mux)
}

func handleIndex() http.HandlerFunc {
	param := page.StatusParams{
		IndexParams: page.IndexParams{
			StylesheetURL: []string{indexStyleURL, statusStyleURL},
		},
		Title:        "Status viewer",
		Calculations: calcs,
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}
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

func fileServerHandler(path string) http.Handler {
	log.Println("serving files at", path)
	// Necessary to strip because file server serves relative to ./public/ folder.
	return http.StripPrefix(path,
		http.FileServer(http.Dir(strings.ReplaceAll(path, "/", ""))))
}

func serveFile(path string) http.HandlerFunc {
	path = strings.TrimPrefix(path, "/")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func loggingHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprintf(w, "You just got %d'ed! :-(", status)
	}
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
