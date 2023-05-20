package page

import (
	"embed"
	"html/template"
	"io"
)

//go:embed html/*
var files embed.FS

// All templates are placed inside the index file using the {{contents}} var name.
const (
	indexFilename = "index.html"
	indexPath     = "html/" + indexFilename
)

// HTML template.
var (
	status = parse("html/status.html")
)

type IndexParams struct {
	StylesheetURL []string
	FaviconURL    string
}

type StatusParams struct {
	IndexParams  IndexParams
	Title        string
}

func Status(w io.Writer, p StatusParams) error {
	return status.Execute(w, p)
}

func parse(fpath string) *template.Template {
	return template.Must(
		template.New(indexFilename).ParseFS(files, indexPath, fpath))
}
