package site

import (
	"embed"
	"html/template"
	"io"
)

//go:embed html/*
var files embed.FS

// All templates are placed inside the index file using the {{content}} var name.
const (
	indexFilename = "index.html"
	indexPath     = "html/" + indexFilename
)

// HTML template.
var (
	status = parse("html/status.html")
)

type StatusPageParam struct {
	Title string
}

func StatusPage(w io.Writer, p StatusPageParam) error {
	return status.Execute(w, p)
}

func parse(fpath string) *template.Template {
	return template.Must(
		template.New(indexFilename).ParseFS(files, indexPath, fpath))
}
