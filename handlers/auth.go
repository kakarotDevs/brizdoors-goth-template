package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
	"golang.org/x/crypto/bcrypt"
)

func isUserLoggedIn(r *http.Request) bool {
	_, ok := internal.GetUserFromSession(r)
	return ok
}

func generateState() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		panic(err) // handle better in production
	}
	return base64.URLEncoding.EncodeToString(b)
}

func HandleRegister(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return Render(w, r, partials.Register(""))
	}

	if r.Method == http.MethodPost {
		ID := uuid.New().String()
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			// Return form with error
			return Render(w, r, partials.Register("Email and password are required"))
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return Render(w, r, partials.Register("Failed to hash password"))
		}

		user := models.User{
			ID:       ID,
			Email:    email,
			Password: string(hashedPassword),
		}

		err = models.CreateUser(&user)
		if err != nil {
			if err == models.ErrUserAlreadyExists {
				return Render(w, r, partials.Register("User with that email already exists"))
			}
			return Render(w, r, partials.Register("Failed to create user"))
		}

		// On success, redirect (HTMX handles redirect headers)
		internal.SetUserSession(w, user.ID)
		w.Header().Set("HX-Redirect", "/lobby")
		w.WriteHeader(http.StatusOK)
		return nil
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func HandleLogin(w http.ResponseWriter, r *http.Request) error {
	email := r.FormValue("email")
	password := r.FormValue("password")

	user, err := models.GetUserByEmail(email)

	if err == models.ErrUserNotFound {
		if r.Header.Get("HX-Request") == "true" {
			return Render(w, r, partials.Register("No user found, create a new account."))
		}

		http.Redirect(w, r, "/register", http.StatusFound)
		return nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return nil
	}

	internal.SetUserSession(w, user.ID)
	http.Redirect(w, r, "/lobby", http.StatusFound)
	return nil
}

func HandleLogout(w http.ResponseWriter, r *http.Request) error {
	internal.ClearUserSession(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
	return nil
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) error {
	// Generate a random state and save it to session for CSRF protection in production
	state := generateState()

	url := auth.GoogleOauthConfig.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
	return nil
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) error {
	ctx := context.Background()

	// Read state and code from query params
	state := r.FormValue("state")
	code := r.FormValue("code")

	// TODO: Verify state matches what you saved in session for CSRF protection
	_ = state // Placeholder to prevent unused variable error

	token, err := auth.GoogleOauthConfig.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	client := auth.GoogleOauthConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return nil
	}
	defer resp.Body.Close()

	var userInfo struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to decode user info: "+err.Error(), http.StatusInternalServerError)
		return nil
	}

	// Store/update user in your database
	models.CreateOrUpdateUser(models.User{
		ID:      userInfo.ID,
		Name:    userInfo.Name,
		Email:   userInfo.Email,
		Picture: userInfo.Picture,
	})

	// Set user session to logged in
	internal.SetUserSession(w, userInfo.ID)

	// Redirect to home or dashboard
	http.Redirect(w, r, "/lobby", http.StatusFound)
	return nil
}

func HandleAuthMenu(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		// Not logged in — return empty or guest UI
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	// Render a minimal partial with user name and links
	return Render(w, r, partials.AuthMenu(user.Name))
}

func HandleAuthMenuToggle(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	// Return the dropdown content
	return Render(w, r, partials.AuthMenuContent(user.Name))
}

func HandleAuthMenuContent(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	user, err := models.GetUserByID(userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	// Return just the dropdown content
	return Render(w, r, partials.AuthMenuContent(user.Name))
}
