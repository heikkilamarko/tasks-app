package ui

import (
	"embed"
	"html/template"
	"io"
	"log/slog"
)

//go:embed templates
var templatesFS embed.FS

type TemplateRenderer struct {
	templates *template.Template
	logger    *slog.Logger
}

func NewTemplateRenderer(logger *slog.Logger) (*TemplateRenderer, error) {
	templates, err := template.New("").
		Funcs(template.FuncMap{
			"dict":          Dict,
			"formattime":    FormatTime,
			"formatisotime": FormatISOTime,
		}).
		ParseFS(templatesFS, "templates/*.html")

	if err != nil {
		return nil, err
	}

	return &TemplateRenderer{templates, logger}, nil
}

func (r *TemplateRenderer) Render(w io.Writer, name string, data any) error {
	if err := r.templates.ExecuteTemplate(w, name, data); err != nil {
		r.logger.Error("render template", "error", err)
		panic(err)
	}
	return nil
}
