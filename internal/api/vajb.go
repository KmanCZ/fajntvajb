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

	areErrors := len(validationErrors) > 0

	validationErrors["Name"] = name
	validationErrors["Description"] = description
	validationErrors["Address"] = address
	validationErrors["Region"] = region
	validationErrors["Date"] = date
	validationErrors["Time"] = timeF
	validationErrors["MinDate"] = time.Now().Format("2006-01-02")

	if areErrors {
		err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

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

	var isJoined bool
	if user != nil {
		isJoined, err = handlers.db.Vajbs.GetIsJoinedToVajb(id, user.ID)
		if err != nil {
			handleWebError(w, err)
			return
		}
	}

	participants, err := handlers.db.Vajbs.GetVajbParticipants(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	for i := 0; i < len(participants); i++ {
		participants[i].ProfilePic.String = files.GetProfilePicPath(participants[i].ProfilePic)
		participants[i].ProfilePic.Valid = true
	}

	err = handlers.tmpl.Render(w, r, "vajb_page", map[string]any{
		"Vajb":            vajb,
		"ImagePath":       files.GetVajbPicPath(vajb.HeaderImage),
		"Date":            vajb.Date.Format("02. 01. 2006 15:04"),
		"Region":          handlers.db.Vajbs.GetFullRegionName(vajb.Region),
		"IsOwner":         isOwner,
		"IsAuthenticated": user != nil,
		"IsJoined":        isJoined,
		"Participants":    participants,
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

func (handlers *handlers) handleEditVajb(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	originalVajb, err := handlers.db.Vajbs.GetVajbByID(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
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
	deleteCurrentImage := r.Form.Get("delete_header_image")

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

	areErrors := len(validationErrors) > 0

	validationErrors["Name"] = name
	validationErrors["Description"] = description
	validationErrors["Address"] = address
	validationErrors["Region"] = region
	validationErrors["Date"] = date
	validationErrors["Time"] = timeF
	validationErrors["MinDate"] = time.Now().Format("2006-01-02")
	validationErrors["Edit"] = true
	validationErrors["ID"] = id
	validationErrors["HasImage"] = originalVajb.HeaderImage.Valid
	validationErrors["ImagePath"] = files.GetVajbPicPath(originalVajb.HeaderImage)

	if areErrors {
		err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	image, _, err := r.FormFile("header_image")
	if err != nil && err != http.ErrMissingFile {
		validationErrors["HeaderImageError"] = "Failed to upload image"
		err = handlers.tmpl.Render(w, r, "vajb_form", validationErrors)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	imgId := originalVajb.HeaderImage.String
	if originalVajb.HeaderImage.Valid && (deleteCurrentImage != "" || image != nil) {
		err = files.DeleteVajbPic(originalVajb.HeaderImage.String)
		if err != nil {
			handleWebError(w, err)
			return
		}
		imgId = ""
	}

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
	err = handlers.db.Vajbs.UpdateVajb(originalVajb.ID, user.ID, name, description, address, region, imgId, finalDate)
	if err != nil {
		handleWebError(w, err)
		return
	}

	http.Redirect(w, r, "/vajb/"+strconv.Itoa(originalVajb.ID), http.StatusSeeOther)
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

	if vajb.HeaderImage.Valid {
		err = files.DeleteVajbPic(vajb.HeaderImage.String)
		if err != nil {
			handleWebError(w, err)
			return
		}
	}

	err = handlers.db.Vajbs.DeleteVajb(id)
	if err != nil {
		handleWebError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/vajb")
}

func (handlers *handlers) handleJoinVajb(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	user := r.Context().Value("user").(*repository.User)
	err = handlers.db.Vajbs.JoinVajb(id, user.ID)
	if err != nil {
		handleWebError(w, err)
		return
	}

	http.Redirect(w, r, "/vajb/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (handlers *handlers) handleUnjoinVajb(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = handlers.tmpl.Render(w, r, "404", nil)
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	user := r.Context().Value("user").(*repository.User)
	err = handlers.db.Vajbs.UnjoinVajb(id, user.ID)
	if err != nil {
		handleWebError(w, err)
		return
	}

	http.Redirect(w, r, "/vajb/"+strconv.Itoa(id), http.StatusSeeOther)
}

func (handlers *handlers) handleVajbExplorePage(w http.ResponseWriter, r *http.Request) {
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
		handleWebError(w, err)
		return
	}

	err = handlers.tmpl.Render(w, r, "explore", map[string]any{
		"Region": region,
		"Name":   name,
		"From":   from,
		"To":     to,
		"Num":    num,
		"Offset": off + number,
		"Vajbs":  vajbs,
	})
	if err != nil {
		handleWebError(w, err)
	}
}
