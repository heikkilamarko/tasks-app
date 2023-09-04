package emailnotifier

import (
	"embed"
	"html/template"
)

//go:embed templates
var TemplatesFS embed.FS

var Templates = template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))
