package handlers


import(
	"net/http"
	"forum/database"
)


func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
    // Ensure the method is POST for creating a post

	
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    // Fetch the session cookie
    cookie, err := r.Cookie("session_id")
    if err != nil {
        http.Error(w, "You must be logged in to create a post", http.StatusUnauthorized)
        return
    }

    // Validate the session
    sessionID := cookie.Value
    userID, loggedIn := sessions[sessionID]
    if !loggedIn {
        http.Error(w, "Invalid session. Please log in again.", http.StatusUnauthorized)
        return
    }

    // Parse the form data
    err = r.ParseForm()
    if err != nil {
        http.Error(w, "Failed to parse form data: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Get the post details
    title := r.FormValue("title")
    content := r.FormValue("content")

    if title == "" || content == "" {
        http.Error(w, "Title and content cannot be empty", http.StatusBadRequest)
        return
    }

    // Insert the post into the database
    _, err = database.DB.Exec(`
        INSERT INTO posts (title, content, user_id) 
        VALUES (?, ?, ?)`, title, content, userID)
    if err != nil {
        http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
        return
    }

    // Redirect back to the home page
    http.Redirect(w, r, "/", http.StatusSeeOther)
}
