package api

import (
	"fajntvajb/internal/database"
	"fajntvajb/internal/files/templates"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/validator"
)

type handlers struct {
	tmpl      *templates.Template
	db        *database.DB
	validator *validator.Validator
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
		tmpl:      templates,
		db:        db,
		validator: validator.New(),
	}
	return &res, nil
}
