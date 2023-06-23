package assets

import "embed"

//go:embed templates/*
var Templates embed.FS

//go:embed static
var StaticFS embed.FS
