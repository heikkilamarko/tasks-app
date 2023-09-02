package internal

import (
	"embed"
	"html/template"
)

//go:embed ui/templates
var UITemplatesFS embed.FS

var UITemplates = template.Must(template.New("").
	Funcs(template.FuncMap{
		"FormatUITime":        FormatUITime,
		"FormatUIDisplayTime": FormatUIDisplayTime,
		"ParseUITime":         ParseUITime,
	}).
	ParseFS(UITemplatesFS, "ui/templates/*.html"),
)
