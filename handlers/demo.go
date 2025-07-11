package handlers

import (
	"net/http"
	"time"

	"github.com/kakarotDevs/brizdoors-goth-template/views/demo"
)

// HandleSpinnerDemo shows the spinner demo page
func HandleSpinnerDemo(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, demo.Index())
}

// HandleDemoLoading simulates a loading action with delay
func HandleDemoLoading(w http.ResponseWriter, r *http.Request) error {
	// Simulate processing time
	time.Sleep(2 * time.Second)

	response := `
		<div class="p-4 bg-green-100 border border-green-400 text-green-700 rounded-lg">
			<h3 class="font-semibold">Success!</h3>
			<p>Data loaded successfully after 2 seconds.</p>
			<p class="text-sm text-green-600 mt-2">Timestamp: ` + time.Now().Format("15:04:05") + `</p>
		</div>
	`

	w.Write([]byte(response))
	return nil
}

// HandleDemoForm processes a form submission with loading state
func HandleDemoForm(w http.ResponseWriter, r *http.Request) error {
	// Simulate processing time
	time.Sleep(1 * time.Second)

	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return nil
	}

	input := r.FormValue("input")
	if input == "" {
		input = "No input provided"
	}

	response := `
		<div class="p-4 bg-blue-100 border border-blue-400 text-blue-700 rounded-lg">
			<h3 class="font-semibold">Form Submitted!</h3>
			<p>You entered: <strong>` + input + `</strong></p>
			<p class="text-sm text-blue-600 mt-2">Processed at: ` + time.Now().Format("15:04:05") + `</p>
		</div>
	`

	w.Write([]byte(response))
	return nil
}

// HandleDemoContent loads content with a delay
func HandleDemoContent(w http.ResponseWriter, r *http.Request) error {
	// Simulate loading time
	time.Sleep(3 * time.Second)

	response := `
		<div class="p-6 bg-bg-light rounded-lg border border-bg-light/30">
			<h3 class="text-xl font-semibold text-text mb-4">Loaded Content</h3>
			<div class="space-y-3 text-text">
				<p>This content was loaded asynchronously using HTMX.</p>
				<p>The branded spinner was shown during the loading process.</p>
				<p class="text-sm text-text-muted">Loaded at: ` + time.Now().Format("15:04:05") + `</p>
			</div>
			<div class="mt-4 flex space-x-2">
				<button
					class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition-colors"
					hx-get="/demo/content"
					hx-target="#content-area"
					hx-indicator="#content-spinner"
				>
					Reload
				</button>
				<button
					class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
					hx-get="/demo/clear"
					hx-target="#content-area"
				>
					Clear
				</button>
			</div>
		</div>
	`

	w.Write([]byte(response))
	return nil
}

// HandleDemoClear clears the content area
func HandleDemoClear(w http.ResponseWriter, r *http.Request) error {
	response := `
		<div class="text-center text-text-muted py-8">
			<p>Content cleared. Click "Load Content" to load new content.</p>
		</div>
	`

	w.Write([]byte(response))
	return nil
}
