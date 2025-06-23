//go:build dev

package main

import (
	"fmt"
	"net/http"
)

func public() http.Handler {
	fmt.Println("Using local public/ directory for static files (dev mode)")
	return http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
}
