package handlers

import (
	"fmt"
	"net/http"
)

// ChatHandler handles chat POST requests from the chatbot UI
func ChatHandler(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return nil
	}

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return nil
	}

	// TODO: Integrate with AI backend here to generate a reply
	reply := generateAIResponse(message)

	// Return chat HTML snippet (user message + AI reply)
	_, err := fmt.Fprintf(w, `
    <div class="self-end bg-orange-200 rounded px-3 py-2 max-w-[70%%] break-words">
      %s
    </div>
    <div class="self-start bg-gray-200 rounded px-3 py-2 max-w-[70%%] break-words mt-1">
      %s
    </div>
  `, message, reply)

	return err
}

// Placeholder function for AI integration
func generateAIResponse(message string) string {
	// TODO: Call your AI service here (e.g., OpenAI, etc.)
	return fmt.Sprintf("This is a placeholder AI response to: %s", message)
}
