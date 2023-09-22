package ui

import (
	"embed"
	"html/template"
)

//go:embed templates
var TemplatesFS embed.FS

var Templates = template.Must(template.New("").
	Funcs(template.FuncMap{
		"RenderEnv":     RenderEnv,
		"RenderTime":    RenderTime,
		"RenderISOTime": RenderISOTime,
	}).
	ParseFS(TemplatesFS, "templates/*.html"),
)
