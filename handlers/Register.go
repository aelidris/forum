package handlers

import (
	"html/template"
	"net/http"
	"regexp"

	"forum/database"
)

var emailRegex = regexp.MustCompile(`^[-!#$%&'*+/0-9=?A-Z^_a-z{|}~](\.?[-!#$%&'*+/0-9=?A-Z^_a-z` + "`" + `{|}~])*@[a-zA-Z0-9](-*\.?[a-zA-Z0-9])*\.[a-zA-Z](-?[a-zA-Z0-9])+$`)

var registerTemplate = template.Must(template.ParseFiles("templates/register.html"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		registerTemplate.Execute(w, nil)
		return
	}

	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Validate input fields
	if email == "" || username == "" || password == "" {
		data := map[string]string{
			"Error":    "All fields are required",
			"Email":    email,
			"Username": username,
		}
		registerTemplate.Execute(w, data)
		return
	}

	// Validate email format
	if !emailRegex.MatchString(email) {
		data := map[string]string{
			"Error":    "Invalid email format. Please try again.",
			"Email":    email,
			"Username": username,
		}
		registerTemplate.Execute(w, data)
		return
	}

	// Hash the password for security (use a proper hashing library like bcrypt)
	hashedPassword := password // Replace this with hashed password using bcrypt

	// Save the user in the database
	_, err = database.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
	if err != nil {
		data := map[string]string{
			"Error":    "Registration failed: " + err.Error(),
			"Email":    email,
			"Username": username,
		}
		registerTemplate.Execute(w, data)
		return
	}

	// Redirect to login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
