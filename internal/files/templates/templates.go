package templates

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func New() (*Template, error) {
	log := logger.Get()
	t, err := template.New("layout.html").ParseFS(files.Files, "templates/layouts/*.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse layout")
		return nil, err
	}

	return &Template{
		templates: t,
	}, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
	log := logger.Get()
	tmpl, err := t.templates.Clone()
	if err != nil {
		log.Error().Err(err).Msg("Failed to clone template")
		return err
	}
	tmpl, err = tmpl.ParseFS(files.Files, "templates/pages/"+name+".html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse template")
		return err
	}
	return tmpl.ExecuteTemplate(w, name+".html", data)

}

