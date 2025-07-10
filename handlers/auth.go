package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kakarotDevs/brizdoors-goth-template/db"
	"github.com/kakarotDevs/brizdoors-goth-template/internal"
	"github.com/kakarotDevs/brizdoors-goth-template/internal/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/models"
	authviews "github.com/kakarotDevs/brizdoors-goth-template/views/auth"
	"github.com/kakarotDevs/brizdoors-goth-template/views/partials"
	"golang.org/x/crypto/bcrypt"
)

// splitName splits a full name into first and last name for Google users
func splitName(full string) (string, string) {
	parts := strings.Fields(full)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}

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
		return Render(w, r, authviews.Register(""))
	}

	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "" || password == "" {
			return Render(w, r, authviews.Register("Email and password are required"))
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return Render(w, r, authviews.Register("Error creating account"))
		}

		hashedPasswordString := string(hashedPassword)

		user := models.User{
			ID:        uuid.New().String(),
			Email:     email,
			Password:  &hashedPasswordString,
			FirstName: "",
			LastName:  "",
			Role:      "user",
			Picture:   "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err = models.CreateUser(r.Context(), db.DB, &user)
		if err != nil {
			if err == models.ErrUserAlreadyExists {
				return Render(w, r, authviews.Register("User with that email already exists"))
			}
			return Render(w, r, authviews.Register("Error creating account"))
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return nil
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return nil
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		msg := "Email and password are required"
		if r.Header.Get("HX-Request") == "true" {
			return Render(w, r, authviews.LoginForm(msg))
		}
		http.Error(w, msg, http.StatusBadRequest)
		return nil
	}

	user, err := models.GetUserByEmail(r.Context(), db.DB, email)
	if err == models.ErrUserNotFound {
		return createAccountAndLogin(w, r, email, password)
	}
	if err != nil {
		return fmt.Errorf("error fetching user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(password)); err != nil {
		msg := "Invalid email or password"
		if r.Header.Get("HX-Request") == "true" {
			return Render(w, r, authviews.LoginForm(msg))
		}
		http.Error(w, msg, http.StatusUnauthorized)
		return nil
	}

	internal.SetUserSession(w, user.ID)
	http.Redirect(w, r, "/lobby", http.StatusFound)
	return nil
}

func HandleLogout(w http.ResponseWriter, r *http.Request) error {
	internal.ClearUserSession(w, r)

	// Set no-cache headers to prevent caching
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")

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

	// Try to find user by email (merge Google and local accounts by email)
	existingUser, err := models.GetUserByEmail(r.Context(), db.DB, userInfo.Email)
	if err == nil && existingUser != nil && (existingUser.GoogleID == "" || existingUser.GoogleID == userInfo.ID) {
		// Upgrade local account: just add GoogleID and update fields
		firstName, lastName := splitName(userInfo.Name)
		existingUser.FirstName = firstName
		existingUser.LastName = lastName
		existingUser.GoogleID = userInfo.ID
		existingUser.Role = "user"
		existingUser.Picture = userInfo.Picture
		existingUser.UpdatedAt = time.Now()
		// Optionally blank password - uncomment next line if you want Google-only login
		existingUser.Password = nil
		err = models.CreateOrUpdateUser(r.Context(), db.DB, existingUser)
		if err != nil {
			log.Printf("Failed to update local user for Google upgrade: %v", err)
			http.Error(w, "Error upgrading account", http.StatusInternalServerError)
			return nil
		}
		internal.SetFlash(w, "Your account has been linked to Google. Please use Google to log in in the future.")
		internal.SetUserSession(w, existingUser.ID)
		http.Redirect(w, r, "/lobby", http.StatusFound)
		return nil
	} else {
		// Find by Google ID first
		user, err := models.GetUserByGoogleID(r.Context(), db.DB, userInfo.ID)
		if err == models.ErrUserNotFound {
			firstName, lastName := splitName(userInfo.Name)

			// New Google-linked user
			googleUser := &models.User{
				ID:        uuid.New().String(),
				GoogleID:  userInfo.ID,
				FirstName: firstName,
				LastName:  lastName,
				Role:      "user",
				Email:     userInfo.Email,
				Picture:   userInfo.Picture,
				Password:  nil,
			}
			err = models.CreateUser(r.Context(), db.DB, googleUser)
			if err != nil {
				log.Printf("Google user create error: %v", err)
				http.Error(w, "Failed to create account", http.StatusInternalServerError)
				return nil
			}
			user = googleUser
		} else if err != nil {
			log.Printf("Google user lookup error: %v", err)
			http.Error(w, "Failed to look up user", http.StatusInternalServerError)
			return nil
		} else {
			// Already registered Google user, update fields
			user.FirstName, user.LastName = splitName(userInfo.Name)
			user.Role = "user"

			user.Email = userInfo.Email
			user.Picture = userInfo.Picture
			err = models.CreateOrUpdateUser(r.Context(), db.DB, user)
			if err != nil {
				log.Printf("Google user update error: %v", err)
				http.Error(w, "Failed to update user", http.StatusInternalServerError)
				return nil
			}
		}

		internal.SetFlash(w, "Your account has been upgraded to Google sign-in!")
		internal.SetUserSession(w, user.ID)
		// Redirect to lobby
		http.Redirect(w, r, "/lobby", http.StatusFound)
		return nil
	}
}

func HandleAuthMenu(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		// Not logged in — return empty or guest UI
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusNoContent)
		return nil
	}

	// Render a minimal partial with user name and links
	return Render(w, r, partials.AuthMenu(user.FirstName, user.Picture, user.Email))
}

func HandleAuthMenuToggle(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	// Return the open dropdown menu with content
	return Render(w, r, partials.AuthMenuOpen(user.FirstName, user.Picture, user.Email))
}

func HandleAuthMenuContent(w http.ResponseWriter, r *http.Request) error {
	userID, ok := internal.GetUserFromSession(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	user, err := models.GetUserByID(r.Context(), db.DB, userID)
	if err != nil {
		internal.ClearUserSession(w, r)
		w.WriteHeader(http.StatusUnauthorized)
		return nil
	}

	// Return just the dropdown content
	return Render(w, r, partials.AuthMenuContent(user.FirstName, user.Picture, user.Email))
}

// createAccountAndLogin creates a new account and logs the user in automatically
func createAccountAndLogin(w http.ResponseWriter, r *http.Request, email, password string) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if r.Header.Get("HX-Request") == "true" {
			return Render(w, r, authviews.Register("Error creating account"))
		}
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return nil
	}

	hashedPasswordString := string(hashedPassword)

	// Create the user
	user := models.User{
		ID:       uuid.New().String(),
		Email:    email,
		Password: &hashedPasswordString,
	}

	err = models.CreateUser(r.Context(), db.DB, &user)
	if err != nil {
		if err == models.ErrUserAlreadyExists {
			// This shouldn't happen since we checked, but handle gracefully
			if r.Header.Get("HX-Request") == "true" {
				return Render(w, r, authviews.Register("Account already exists. Please try logging in."))
			}
			http.Error(w, "Account already exists", http.StatusConflict)
			return nil
		}
		if r.Header.Get("HX-Request") == "true" {
			return Render(w, r, authviews.Register("Error creating account"))
		}
		http.Error(w, "Error creating account", http.StatusInternalServerError)
		return nil
	}

	// Automatically log the user in
	internal.SetUserSession(w, user.ID)

	// Redirect to lobby
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/lobby")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, "/lobby", http.StatusFound)
	}
	return nil
}
