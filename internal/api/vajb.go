package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/repository"
	"fajntvajb/internal/validator"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func (handlers *handlers) handleNewVajbPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, r, "vajb_form", map[string]any{
		"MinDate": time.Now().Format("2006-01-02"),
	})
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleNewVajb(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		handleWebError(w, err)
		return
	}

	name := r.Form.Get("name")
	description := r.Form.Get("description")
	address := r.Form.Get("address")
	region := r.Form.Get("region")
	date := r.Form.Get("date")
	timeF := r.Form.Get("time")

	err = handlers.validator.ValidateVajb(&validator.Vajb{
		Name:        name,
		Description: description,
		Address:     address,
		Region:      region,
		Date:        date,
		Time:        timeF,
	})
	validationErrors := make(map[string]any)
	if err != nil {
		validationErrors = handlers.validator.HandleVajbValidationError(err)
	}

	if len(validationErrors) > 0 {
		validationErrors["Name"] = name
		validationErrors["Description"] = description
		validationErrors["Address"] = address
		validationErrors["Region"] = region
		validationErrors["Date"] = date
		validationErrors["Time"] = timeF
		validationErrors["MinDate"] = time.Now().Format("2006-01-02")
		err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	validationErrors["Name"] = name
	validationErrors["Description"] = description
	validationErrors["Address"] = address
	validationErrors["Region"] = region
	validationErrors["Date"] = date
	validationErrors["Time"] = timeF
	validationErrors["MinDate"] = time.Now().Format("2006-01-02")

	image, _, err := r.FormFile("header_image")
	if err != nil && err != http.ErrMissingFile {
		validationErrors["HeaderImageError"] = "Failed to upload image"
		err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	var imgId string
	if image != nil {
		defer image.Close()
		imageBytes, err := io.ReadAll(image)
		if err != nil {
			validationErrors["HeaderImageError"] = "Failed to upload image"
			err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
			if err != nil {
				handleWebError(w, err)
			}
			return
		}
		imgId = uuid.New().String()
		err = files.UploadVajbPic(imgId, imageBytes)
		if err != nil {
			validationErrors["HeaderImageError"] = "Failed to upload image"
			err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
			if err != nil {
				handleWebError(w, err)
			}
			return
		}
	}

	user := r.Context().Value("user").(*repository.User)
	finalDate, err := time.Parse("2006-01-02 15:04", date+" "+timeF)
	if err != nil {
		handleWebError(w, err)
		return
	}
	vajb, err := handlers.db.Vajbs.CreateVajb(user.ID, name, description, address, region, imgId, finalDate)
	if err != nil {
		handleWebError(w, err)
		return
	}

	http.Redirect(w, r, "/vajb/"+strconv.Itoa(vajb.ID), http.StatusSeeOther)
}

func (handlers *handlers) handleVajbPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	vajb, err := handlers.db.Vajbs.GetVajbByID(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	if vajb == nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	var user *repository.User
	if r.Context().Value("user") != nil {
		user = r.Context().Value("user").(*repository.User)
	}

	var isOwner bool
	if user != nil {
		isOwner = user.ID == vajb.CreatorID
	}

	err = handlers.tmpl.Render(w, r, "vajb_page", map[string]any{
		"Vajb":            vajb,
		"ImagePath":       files.GetVajbPicPath(vajb.HeaderImage),
		"Date":            vajb.Date.Format("02. 01. 2006 15:04"),
		"Region":          handlers.db.Vajbs.GetFullRegionName(vajb.Region),
		"IsOwner":         isOwner,
		"IsAuthenticated": user != nil,
	})
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleVajbEditPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	vajb, err := handlers.db.Vajbs.GetVajbByID(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	if vajb == nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	var user *repository.User
	if r.Context().Value("user") != nil {
		user = r.Context().Value("user").(*repository.User)
	}

	var isOwner bool
	if user != nil {
		isOwner = user.ID == vajb.CreatorID
	}

	if !isOwner {
		http.Redirect(w, r, "/vajb/"+strconv.Itoa(vajb.ID), http.StatusSeeOther)
		return
	}

	err = handlers.tmpl.Render(w, r, "vajb_form", map[string]any{
		"ID":          vajb.ID,
		"Name":        vajb.Name,
		"Description": vajb.Description,
		"Address":     vajb.Address,
		"Region":      vajb.Region,
		"Date":        vajb.Date.Format("2006-01-02"),
		"Time":        vajb.Date.Format("15:04"),
		"MinDate":     time.Now().Format("2006-01-02"),
		"ImagePath":   files.GetVajbPicPath(vajb.HeaderImage),
		"HasImage":    vajb.HeaderImage.Valid,
		"Edit":        true,
	})
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleDeleteVajb(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	vajb, err := handlers.db.Vajbs.GetVajbByID(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	if vajb == nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	user := r.Context().Value("user").(*repository.User)
	if user.ID != vajb.CreatorID {
		handleWebError(w, fmt.Errorf("user is not the creator of the vajb"))
		return
	}

	err = handlers.db.Vajbs.DeleteVajb(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/vajb")
}
