package api

import (
	"net/http"
	"os"

	"fajntvajb/internal/database"
	"fajntvajb/internal/files/templates"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/validator"

	"github.com/gorilla/sessions"
)

type handlers struct {
	tmpl      *templates.Template
	db        *database.DB
	validator *validator.Validator
	session   *sessions.CookieStore
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

	session := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session.Options.HttpOnly = true
	session.Options.Secure = true
	session.Options.SameSite = http.SameSiteLaxMode

	validator, err := validator.New()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create validator")
		return nil, err
	}

	res := handlers{
		tmpl:      templates,
		db:        db,
		validator: validator,
		session:   session,
	}
	return &res, nil
}
