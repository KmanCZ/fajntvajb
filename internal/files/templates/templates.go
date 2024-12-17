package templates

import (
	"fajntvajb/internal/files"
	"html/template"
	"io"
)

type Template struct {
	templates *template.Template
}

func New() (*Template, error) {
	t, err := template.New("layout.html").ParseFS(files.Files, "templates/layouts/*.html")
	if err != nil {
		return nil, err
	}

	return &Template{
		templates: t,
	}, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}) error {
    tmpl, err := t.templates.Clone()
	if err != nil {
		return err
	}
    tmpl, err = tmpl.ParseFS(files.Files, "templates/pages/"+name+".html")
	if err != nil {
		return err
	}
    return tmpl.ExecuteTemplate(w, name+".html", data)

}