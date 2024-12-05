package handlers


import(
    "net/http"
    "forum/database"
)


func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    
    if r.Method == http.MethodGet {
        http.ServeFile(w, r, "templates/register.html") // Render registration form
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
        http.Error(w, "All fields are required", http.StatusBadRequest)
        return
    }

    // Hash the password for security (use a proper hashing library like bcrypt)
    hashedPassword := password // Replace this with hashed password using bcrypt

    // Save the user in the database
    _, err = database.DB.Exec("INSERT INTO users (email, username, password) VALUES (?, ?, ?)", email, username, hashedPassword)
    if err != nil {
        http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Redirect to login page
    http.Redirect(w, r, "/", http.StatusSeeOther)

}
