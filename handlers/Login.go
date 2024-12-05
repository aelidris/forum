package handlers

import (

	"net/http"
    "github.com/gofrs/uuid"
	"time"
	"forum/database"

)

var sessions = make(map[string]int) // session_id -> user_id mapping

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        http.ServeFile(w, r, "templates/login.html") // Render login form
        return
    }

    // Parse form data
    err := r.ParseForm()
    if err != nil {
        http.Error(w, "Invalid form data", http.StatusBadRequest)
        return
    }

    username := r.FormValue("username")
    password := r.FormValue("password")

    // Validate input
    if username == "" || password == "" {
        http.Error(w, "Both username and password are required", http.StatusBadRequest)
        return
    } 

    // Fetch the user from the database
    var userID int
    var storedPassword string
    err = database.DB.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userID, &storedPassword)
    if err != nil {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Validate the password (for now, compare plaintext passwords; use bcrypt in production)
    if password != storedPassword {
        http.Error(w, "Invalid username or password", http.StatusUnauthorized)
        return
    }

    // Generate a new session ID
    sessionID, err := uuid.NewV4()
    if err != nil {
        http.Error(w, "Failed to create session", http.StatusInternalServerError)
        return
    }

    // Store the session ID and associated user ID
    sessions[sessionID.String()] = userID

    // Set the session cookie
    http.SetCookie(w, &http.Cookie{
        Name:    "session_id",
        Value:   sessionID.String(),
        Expires: time.Now().Add(24 * time.Hour), // Session valid for 24 hours
    })

    // Redirect to the home page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
