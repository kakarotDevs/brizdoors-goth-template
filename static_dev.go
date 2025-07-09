//go:build dev

package main

import (
	"fmt"
	"net/http"
)

func public() http.Handler {
	fmt.Println("Using local public/ directory for static files (dev mode)")
	
	// Create a custom handler that adds no-cache headers for CSS files in dev mode
	fs := http.FileServer(http.Dir("public"))
	return http.StripPrefix("/public/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add no-cache headers for CSS files to enable hot reloading
		if r.URL.Path == "/styles.css" {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
		}
		fs.ServeHTTP(w, r)
	}))
}
