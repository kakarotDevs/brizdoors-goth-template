package handlers

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/views/chat"
)

// HandleChatPage renders the main chat interface
func HandleChatPage(w http.ResponseWriter, r *http.Request) error {
	return Render(w, r, chat.Index())
}
