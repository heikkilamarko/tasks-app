package internal

import (
	"embed"
	"html/template"
)

//go:embed ui/templates
var UITemplatesFS embed.FS

var UITemplates = template.Must(template.ParseFS(UITemplatesFS, "ui/templates/*.html"))
