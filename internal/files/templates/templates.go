package templates

import (
	"html/template"
	"net/http"

	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/repository"
)

type Template struct {
	templates  *template.Template
	components *template.Template
}

func New() (*Template, error) {
	log := logger.Get()
	t, err := template.New("layout.html").ParseFS(files.Files, "templates/layouts/*.html", "templates/components/*.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse layout")
		return nil, err
	}

	components, err := template.ParseFS(files.Files, "templates/components/*.html")
	if err != nil {
		log.Error().Err(err).Msg("Failed to parse components")
		return nil, err
	}

	return &Template{
		templates:  t,
		components: components,
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

		data["ProfilePicPath"] = files.GetProfilePicPath(user.ProfilePic)
	}

	return tmpl.ExecuteTemplate(w, name+".html", data)
}

func (t *Template) RenderComponent(w http.ResponseWriter, name string, data map[string]any) error {
	return t.components.ExecuteTemplate(w, name+".html", data)
}
