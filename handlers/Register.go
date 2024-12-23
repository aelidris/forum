package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"forum/database"
)

var (
	emailRegex    = regexp.MustCompile(`^[-!#$%&'*+/0-9=?A-Z^_a-z{|}~](\.?[-!#$%&'*+/0-9=?A-Z^_a-z` + "`" + `{|}~])*@[a-zA-Z0-9](-*\.?[a-zA-Z0-9])*\.[a-zA-Z](-?[a-zA-Z0-9])+$`)
	usernameRegex = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]{7,29}$"`) // Username regex
)

var registerTemplate = template.Must(template.ParseFiles("templates/register.html"))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Function to generate dynamic script
	getAlertScript := func(message string) string {
		return fmt.Sprintf(`<script>alert('%s')</script>`, message)
	}
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

	LocalPart, DomainPart := "", ""
	if strings.Contains(email, "@") {
		parts := strings.Split(email, "@")

		// Check if exactly one "@" symbol exists
		if len(parts) == 2 {
			LocalPart = parts[0]
			DomainPart = parts[1]
		}
	}

	// Validate email format
	if !emailRegex.MatchString(email) || len(email) > 320 || len(LocalPart) > 64 || len(DomainPart) > 255 {
		script := getAlertScript("Invalid email format. Please try again.")

		fmt.Fprint(w, script)

		data := map[string]string{
			"Error":    "Invalid email format. Please try again.",
			"Email":    email,
			"Username": username,
		}
		registerTemplate.Execute(w, data)

		return
	}

	// Validate username format
	if !usernameRegex.MatchString(username) {
		script := getAlertScript("Invalid username format")
		fmt.Fprint(w, script)
		data := map[string]string{
			"Error":    "Invalid username.",
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
