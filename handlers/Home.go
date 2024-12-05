package handlers

import(
	"net/http"
	"html/template"
	_ "github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"	
	"forum/database"
)



type Post struct {
	ID  int
	Title   string
	Content string
	Username string
	CreatedAt string
}

 


func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Fetch session cookie
    cookie, err := r.Cookie("session_id")
    var loggedIn bool
    var userID int

    if err == nil {
        sessionID := cookie.Value
        userID, loggedIn = sessions[sessionID] // Check if the session is valid
    }

    // Fetch posts from the database
    rows, err := database.DB.Query(`
        SELECT posts.id, posts.title, posts.content, users.username, posts.created_at
        FROM posts
        JOIN users ON posts.user_id = users.id
        ORDER BY posts.created_at DESC
    `)
    if err != nil {
        http.Error(w, "Failed to fetch posts: "+err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var posts []Post
    for rows.Next() {
        var post Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.Username, &post.CreatedAt)
        if err != nil {
            http.Error(w, "Error reading posts: "+err.Error(), http.StatusInternalServerError)
            return
        }
        posts = append(posts, post)
    }

    // Check if no posts exist
    if len(posts) == 0 {
        posts = nil // Explicitly set to nil if no posts are available
    }

    // Render the home page with different options for logged-in and guest users
    tpl, err := template.ParseFiles("templates/home.html")
    if err != nil {
        http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
        return
    }

	data := struct {
        LoggedIn bool
        UserID   int // Include user ID in the template data
        Posts    []Post
    }{
        LoggedIn: loggedIn,
        UserID:   userID,
        Posts:    posts,
    }

    tpl.Execute(w, data)
    
}
