package handlers

import (
	"fmt"
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/chat"
)

// HandleChatPage renders the main chat interface
func HandleChatPage(w http.ResponseWriter, r *http.Request) error {
	isDarkMode := GetThemeFromRequest(r)
	return Render(w, r, chat.Index(isDarkMode))
}

// HandleChatMessage handles individual chat messages
func HandleChatMessage(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return nil
	}

	message := r.FormValue("message")
	if message == "" {
		http.Error(w, "Empty message", http.StatusBadRequest)
		return nil
	}

	// Generate AI response
	reply := generateAIResponse(message)

	// Return ChatGPT-style message HTML
	_, err := fmt.Fprintf(w, `
		<!-- User Message -->
		<div class="px-4 py-6">
			<div class="flex gap-4 max-w-3xl mx-auto">
				<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center flex-shrink-0">
					<svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
					</svg>
				</div>
				<div class="flex-1">
					<div class="text-text whitespace-pre-wrap">%s</div>
				</div>
			</div>
		</div>

		<!-- AI Response -->
		<div class="px-4 py-6 bg-bg-light/30">
			<div class="flex gap-4 max-w-3xl mx-auto">
				<div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full flex items-center justify-center flex-shrink-0">
					<svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 24 24">
						<path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
					</svg>
				</div>
				<div class="flex-1">
					<div class="text-text whitespace-pre-wrap leading-relaxed">%s</div>
				</div>
			</div>
		</div>
	`, message, reply)

	return err
}

// HandleNewChat clears the chat and starts fresh
func HandleNewChat(w http.ResponseWriter, r *http.Request) error {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return nil
	}

	// Return empty chat - messages will appear as user starts chatting
	_, err := fmt.Fprintf(w, `<!-- Chat cleared - ready for new conversation -->`)

	return err
}

// ChatHandler - keeping for backward compatibility
func ChatHandler(w http.ResponseWriter, r *http.Request) error {
	return HandleChatMessage(w, r)
}

// Placeholder function for AI integration
func generateAIResponse(message string) string {
	// TODO: Call your AI service here (e.g., OpenAI, etc.)
	return fmt.Sprintf("This is a placeholder AI response to: %s", message)
}
