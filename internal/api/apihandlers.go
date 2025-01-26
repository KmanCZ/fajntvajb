package api

import (
	"net/http"
	"strconv"
	"time"

	"fajntvajb/internal/logger"
)

func (handlers *handlers) handleVajbExplore(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	name := r.URL.Query().Get("name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")
	num := r.URL.Query().Get("num")
	offset := r.URL.Query().Get("offset")

	toDate, err := time.Parse("2006-01-02", to)
	if err != nil {
		toDate = time.Time{}
	}
	fromDate, err := time.Parse("2006-01-02", from)
	if err != nil {
		fromDate = time.Time{}
	}
	number, err := strconv.Atoi(num)
	if err != nil {
		number = 9
	}
	off, err := strconv.Atoi(offset)
	if err != nil {
		off = 0
	}

	vajbs, err := handlers.db.Vajbs.GetVajbs(name, region, fromDate, toDate, number, off)
	if err != nil {
		handleAPIError(w, err)
		return
	}

	err = handlers.tmpl.RenderComponent(w, "vajb_list", map[string]any{
		"Region": region,
		"Name":   name,
		"From":   from,
		"To":     to,
		"Num":    num,
		"Offset": off + number,
		"Vajbs":  vajbs,
	})
	if err != nil {
		handleAPIError(w, err)
	}
}

func handleAPIError(w http.ResponseWriter, err error) {
	log := logger.Get()
	log.Error().Err(err).Msg("Failed to render page")
	http.Error(w, "Something went wrong", http.StatusInternalServerError)
}
