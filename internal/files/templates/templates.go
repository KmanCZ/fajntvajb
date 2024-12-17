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

func (t *Template) Render(w io.Writer, name string, data interface{}) {
	//TODO: handle errors
    tmpl := template.Must(t.templates.Clone())
    tmpl = template.Must(tmpl.ParseFS(files.Files, "templates/pages/"+name+".html"))
    tmpl.ExecuteTemplate(w, name+".html", data)
}