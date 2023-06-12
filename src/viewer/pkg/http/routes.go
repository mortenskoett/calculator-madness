package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strings"
	"viewer/pkg/page"
)

/* Handling of all HTTP end points. */

const (
	indexStyleURL  string = "/public/indexstyle.css"
	statusStyleURL string = "/public/statusstyle.css"
	faviconURL     string = "/public/images/calculator-crop.ico"
)

func (s *server) attachRoutes() {
	s.mux.HandleFunc("/", s.logHandlerFunc(s.handleIndex()))
	s.mux.HandleFunc("/favicon.ico", s.logHandlerFunc(s.serveFile(faviconURL)))
	s.mux.Handle("/public/", s.logHandler(s.fileServerHandler("/public/")))
	s.mux.HandleFunc("/ws", s.wsmanager.HandleWS())
}

func (s *server) handleIndex() http.HandlerFunc {
	param := page.StatusParams{
		IndexParams: page.IndexParams{
			FaviconURL:    faviconURL,
			StylesheetURL: []string{indexStyleURL, statusStyleURL},
		},
		Title: "Calculator Web Viewer",
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			s.errorHandler(w, r, http.StatusNotFound)
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

func (s *server) fileServerHandler(path string) http.Handler {
	log.Println("serving files at", path)
	// Necessary to strip because file server serves relative to ./public/ folder.
	return http.StripPrefix(path,
		http.FileServer(http.Dir(strings.ReplaceAll(path, "/", ""))))
}

func (s *server) serveFile(path string) http.HandlerFunc {
	path = strings.TrimPrefix(path, "/")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path)
	})
}

func (s *server) logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func (s *server) logHandlerFunc(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}

func (s *server) errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprintf(w, "You just got %d'ed! :-(", status)
	}
}
