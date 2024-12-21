package api

import (
	"fajntvajb/internal/database"
	"fajntvajb/internal/files"
	"fajntvajb/internal/files/templates"
	"fajntvajb/internal/logger"
	"net/http"
)

type handlers struct {
	tmpl *templates.Template
	db   *database.DB
}

func NewHandlers() (*handlers, error) {
	log := logger.Get()
	templates, err := templates.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create templates")
		return nil, err
	}

	db, err := database.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create database")
		return nil, err
	}

	res := handlers{
		tmpl: templates,
		db:   db,
	}
	return &res, nil
}

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "index", nil)
	if err != nil {
		handleError(w, err)
	}
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {

	rows, err := handlers.db.GetRows()
	if err != nil {
		handleError(w, err)
		return
	}
	err = handlers.tmpl.Render(w, "auth", struct {
		Name string
		Auth bool
		Rows []string
	}{
		Name: "john doe",
		Auth: true,
		Rows: rows,
	})

	if err != nil {
		handleError(w, err)
	}
}

func (handlers *handlers) handleHTMXPostTest(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}
	value := r.FormValue("name")
	err = handlers.db.InsertRow(value)
	if err != nil {
		return
	}
	w.Write([]byte("<li>" + value + "</li>"))
}

func (handlers *handlers) handleHTMXTest(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("test"))
	if err != nil {
		handleError(w, err)
	}
}

func handleError(w http.ResponseWriter, err error) {
	log := logger.Get()
	log.Error().Err(err).Msg("Failed to render page")
	file, err := files.Files.ReadFile("templates/pages/error.html")
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	_, err = w.Write(file)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
}
