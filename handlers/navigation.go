package handlers

import (
	"fmt"
	"net/http"
)

// HandleMobileNav returns the mobile navigation menu content
func HandleMobileNav(w http.ResponseWriter, r *http.Request) error {
	// Simple mobile menu that appears when requested
	_, err := fmt.Fprintf(w, `
		<nav class="px-4 py-4 space-y-3 transition-all duration-200">
			<a href="/contact" class="block text-text-muted hover:text-text transition-colors duration-200 font-medium py-2">Contact</a>
			<button
				class="block w-full text-text-muted hover:text-text transition-colors duration-200 font-medium py-2 text-left"
				hx-get="/nav/mobile-close"
				hx-target="#mobile-menu"
				hx-swap="innerHTML"
			>Close Menu</button>
		</nav>
	`)
	return err
}

// HandleMobileNavClose returns empty content to close the mobile menu
func HandleMobileNavClose(w http.ResponseWriter, r *http.Request) error {
	// Return empty content to hide the mobile menu
	_, err := fmt.Fprintf(w, "")
	return err
}
