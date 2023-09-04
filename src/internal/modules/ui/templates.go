package ui

import (
	"embed"
	"html/template"
)

//go:embed templates
var TemplatesFS embed.FS

var Templates = template.Must(template.New("").
	Funcs(template.FuncMap{
		"FormatUITime":        FormatUITime,
		"FormatUIDisplayTime": FormatUIDisplayTime,
		"ParseUITime":         ParseUITime,
	}).
	ParseFS(TemplatesFS, "templates/*.html"),
)
