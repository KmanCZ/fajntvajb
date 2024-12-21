package files

import "embed"

//go:embed static/* templates/*
var Files embed.FS

//go:embed migrations/*
var Migrations embed.FS
