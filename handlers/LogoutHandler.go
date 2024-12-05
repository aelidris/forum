package handlers


import (
	"net/http"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Remove the session cookie by setting it with an expired date
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0), // Set to a past time
		HttpOnly: true,
	})

	// Redirect the user to the home page or login page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}