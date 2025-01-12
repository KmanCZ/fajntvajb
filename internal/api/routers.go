package api

import (
	"fajntvajb/internal/files"
	"fajntvajb/internal/logger"
	"io/fs"
	"net/http"
)

func New() (http.Handler, error) {
	log := logger.Get()

	r := http.NewServeMux()
	handlers, err := NewHandlers()
	if err != nil {
		log.Error().Err(err).Msg("Failed to create handlers")
		return nil, err
	}

	static, err := fs.Sub(files.Files, "static")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create static file server")
		return nil, err
	}

	// Serve static files from /static
	r.Handle("GET /static/", http.StripPrefix("/static", http.FileServerFS(static)))

	// Define page routes
	r.HandleFunc("GET /auth", handlers.requireAuthMiddleware(handlers.handleAuthPage))
	r.HandleFunc("GET /auth/register", handlers.requireNoAuthMiddleware(handlers.handleRegisterPage))
	r.HandleFunc("POST /auth/register", handlers.requireNoAuthMiddleware(handlers.handleRegister))
	r.HandleFunc("GET /auth/login", handlers.requireNoAuthMiddleware(handlers.handleLoginPage))
	r.HandleFunc("POST /auth/login", handlers.requireNoAuthMiddleware(handlers.handleLogin))
	r.HandleFunc("DELETE /auth/logout", handlers.requireAuthMiddleware(handlers.handleLogout))
	r.HandleFunc("GET /auth/profile", handlers.requireAuthMiddleware(handlers.handleProfilePage))
	r.HandleFunc("POST /auth/profile/displayname", handlers.requireAuthMiddleware(handlers.handleDisplayNameEdit))
	r.HandleFunc("POST /auth/profile/password", handlers.requireAuthMiddleware(handlers.handlePasswordEdit))
	r.HandleFunc("POST /auth/profile/profilepicture", handlers.requireAuthMiddleware(handlers.handleProfilePictureEdit))
	r.HandleFunc("POST /auth/profile/delete", handlers.requireAuthMiddleware(handlers.handleDeleteAccount))

	r.HandleFunc("GET /vajb/new", handlers.requireAuthMiddleware(handlers.handleNewVajbPage))
	r.HandleFunc("POST /vajb/new", handlers.requireAuthMiddleware(handlers.handleNewVajb))
	r.HandleFunc("GET /vajb/{id}", handlers.handleVajbPage)
	r.HandleFunc("DELETE /vajb/{id}", handlers.requireAuthMiddleware(handlers.handleDeleteVajb))
	r.HandleFunc("GET /vajb/{id}/edit", handlers.requireAuthMiddleware(handlers.handleVajbEditPage))
	r.HandleFunc("POST /vajb/{id}/edit", handlers.requireAuthMiddleware(handlers.handleEditVajb))
	r.HandleFunc("GET /vajb/{id}/join", handlers.requireAuthMiddleware(handlers.handleJoinVajb))

	r.HandleFunc("GET /", handlers.handleLandingPage)

	// Define API routes
	r.HandleFunc("GET /api/test", handlers.handleHTMXTest)

	return handlers.authenticateMiddleware(r), nil
}
