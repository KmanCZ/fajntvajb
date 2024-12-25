package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/repository"
	"fajntvajb/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, r, "index", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, r, "register", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleWebError(w, err)
		return
	}

	username := r.Form.Get("username")
	displayName := r.Form.Get("display_name")
	password := r.Form.Get("password")
	passwordConfirmation := r.Form.Get("password2")

	err = handlers.validator.ValidateUser(&validator.User{
		Username:    username,
		DisplayName: displayName,
		Password:    password,
	})

	validationErrors := make(map[string]any)
	if err != nil {
		validationErrors = handlers.validator.HandleUserValidationError(err)
	}
	if password != passwordConfirmation {
		validationErrors["Password2Error"] = "Passwords do not match"
	}

	user, err := handlers.db.Users.GetUserByUsername(username)
	if err != nil {
		handleWebError(w, err)
		return
	}
	if user != nil {
		validationErrors["UsernameError"] = "Username is already taken"
	}

	if len(validationErrors) > 0 {
		// Re-render the register page with the validation errors and the user's input
		validationErrors["Username"] = username
		validationErrors["DisplayName"] = displayName

		err = handlers.tmpl.Render(w, r, "register", validationErrors)

		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	_, err = handlers.db.Users.CreateUser(username, displayName, password)
	if err != nil {
		handleWebError(w, err)
		return
	}
	http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
}

func (handlers *handlers) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, r, "login", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleWebError(w, err)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	user, err := handlers.db.Users.GetUserByUsername(username)
	if err != nil {
		handleWebError(w, err)
		return
	}
	if user == nil {
		// Re-render the login page with an error message
		err = handlers.tmpl.Render(w, r, "login", map[string]any{
			"UsernameError": "Username does not exist",
			"Username":      username,
		})
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Re-render the login page with an error message
		err = handlers.tmpl.Render(w, r, "login", map[string]any{
			"PasswordError": "Incorrect password",
			"Username":      username,
		})
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	session, _ := handlers.session.Get(r, "session")
	session.Values["userId"] = user.ID
	err = session.Save(r, w)
	if err != nil {
		handleWebError(w, err)
		return
	}

	http.Redirect(w, r, "/auth", http.StatusSeeOther)
}

func (handlers *handlers) handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := handlers.session.Get(r, "session")
	delete(session.Values, "userId")
	err := session.Save(r, w)
	if err != nil {
		handleWebError(w, err)
		return
	}

	w.Header().Set("HX-Redirect", "/login")
}

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, r, "auth", map[string]any{
		"Name": "test",
	})

	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleProfilePage(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(*repository.User)
	err := handlers.tmpl.Render(w, r, "profile", map[string]any{
		"DisplayName": user.DisplayName,
	})
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleDisplayNameEdit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		handleWebError(w, err)
		return
	}

	displayName := r.Form.Get("display_name")
	password := r.Form.Get("password")
	user := r.Context().Value("user").(*repository.User)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Re-render the profile page with an error message
		err = handlers.tmpl.Render(w, r, "profile", map[string]any{
			"DisplayName":   displayName,
			"PasswordError": "Incorrect password",
		})
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	if displayName == user.DisplayName {
		// Re-render the profile page with an error message
		err = handlers.tmpl.Render(w, r, "profile", map[string]any{
			"DisplayName":      displayName,
			"DisplayNameError": "Display name is the same",
		})
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	err = handlers.validator.ValidateUsername(displayName)
	if err != nil {
		// Re-render the profile page with an error message
		err = handlers.tmpl.Render(w, r, "profile", map[string]any{
			"DisplayName":      displayName,
			"DisplayNameError": err.Error(),
		})
		if err != nil {
			handleWebError(w, err)
		}
		return
	}

	err = handlers.db.Users.UpdateDisplayName(user.ID, displayName)
	if err != nil {
		handleWebError(w, err)
		return
	}

	user.DisplayName = displayName

	err = handlers.tmpl.Render(w, r, "profile", map[string]any{
		"DisplayName": displayName,
		"Success":     "Display name updated",
	})
	if err != nil {
		handleWebError(w, err)
	}
}

func handleWebError(w http.ResponseWriter, err error) {
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
	}
}
