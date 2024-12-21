package main

import (
	"fmt"
	"net/http"

	"forum/database"

	"forum/handlers"
)

func main() {
	database.InitDB()

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/make-post", handlers.CreatePostHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("server starting at: http://localhost:8080/")
	http.ListenAndServe(":8080", nil)
}
