package auth

import (
	"net/http"

	"github.com/kakarotDevs/brizdoors-goth-template/internal"
)

// RequireAuth ensures the user is authenticated before accessing the route
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok := internal.GetUserFromSession(r)
		if !ok {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
