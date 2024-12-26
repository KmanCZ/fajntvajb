package files

import (
	"database/sql"
	"embed"
)

//go:embed static/* templates/*
var Files embed.FS

//go:embed migrations/*
var Migrations embed.FS

func GetProfilePicPath(profilePicName sql.NullString) string {
	if profilePicName.Valid {
		return "/pictures/" + profilePicName.String
	}
	return "/static/img/blank-profile-picture.png"
}
