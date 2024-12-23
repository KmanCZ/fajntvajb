package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"fajntvajb/internal/validator"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func (handlers *handlers) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "index", nil)
	if err != nil {
		handleWebError(w, err)
	}
}

func (handlers *handlers) handleRegisterPage(w http.ResponseWriter, _ *http.Request) {
	err := handlers.tmpl.Render(w, "register", nil)
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

	validationErrors := make(map[string]string)
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

		err = handlers.tmpl.Render(w, "register", validationErrors)

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

func (handlers *handlers) handleLoginPage(w http.ResponseWriter, _ *http.Request) {
	err := handlers.tmpl.Render(w, "login", nil)
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
		err = handlers.tmpl.Render(w, "login", map[string]string{
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
		err = handlers.tmpl.Render(w, "login", map[string]string{
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

func (handlers *handlers) handleAuthPage(w http.ResponseWriter, r *http.Request) {
	err := handlers.tmpl.Render(w, "auth", struct {
		Name string
		Auth bool
	}{
		Name: "john doe",
		Auth: true,
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
