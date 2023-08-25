package internal

import (
	"embed"
	"html/template"
)

//go:embed ui/templates
var TemplatesFS embed.FS

//go:embed ui/static
var StaticFS embed.FS

var UITemplates = template.Must(template.ParseFS(TemplatesFS, "ui/templates/*.html"))
