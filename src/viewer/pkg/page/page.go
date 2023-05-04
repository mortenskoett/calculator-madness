package page

import (
	"embed"
	"html/template"
	"io"
)

//go:embed html/*
var files embed.FS

// All templates are placed inside the index file using the {{content}} var name.
const (
	indexFilename = "index.gohtml"
	indexPath     = "html/" + indexFilename
)

// HTML template.
var (
	status = parse("html/status.gohtml")
)

type IndexParams struct {
	StylesheetURL []string
	FaviconURL    string
}

type Progress struct {
	Current int
	Outof   int
}

type Calculation struct {
	CalculationID string
	MessageID     string
	CreatedTime   string
	Equation      string
	Progress      Progress
	Result        string
}

type StatusParams struct {
	IndexParams  IndexParams
	Title        string
	Calculations []Calculation
}

func Status(w io.Writer, p StatusParams) error {
	return status.Execute(w, p)
}

func parse(fpath string) *template.Template {
	return template.Must(
		template.New(indexFilename).ParseFS(files, indexPath, fpath))
}
