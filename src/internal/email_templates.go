package internal

import (
	"embed"
	"html/template"
)

//go:embed email/templates
var EmailTemplatesFS embed.FS

var EmailTemplates = template.Must(template.ParseFS(EmailTemplatesFS, "email/templates/*.html"))
