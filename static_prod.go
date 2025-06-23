//go:build !dev
// +build !dev

package main

import (
	"embed"
	"net/http"
)

//go:embed public/*
var publicFS embed.FS

func public() http.Handler {
	fs := http.FileServer(http.FS(publicFS))
	return http.StripPrefix("/public/", fs)
}
