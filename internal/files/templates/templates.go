package templates

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/repository"
	"html/template"
	"net/http"
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

func (t *Template) Render(w http.ResponseWriter, r *http.Request, name string, data map[string]any) error {
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

	if data == nil {
		data = make(map[string]any)
	}

	user, ok := r.Context().Value("user").(*repository.User)
	if !ok {
		data["Auth"] = false
	} else {
		data["Auth"] = true
		data["User"] = user
	}

	return tmpl.ExecuteTemplate(w, name+".html", data)

}
