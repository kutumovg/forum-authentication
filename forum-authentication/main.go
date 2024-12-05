package main

import (
	"log"
	"net/http"

	"forum-authentication/handlers"
	"forum-authentication/models"
)

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	models.SetDB(db)

	// Routes
	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/create_post", handlers.CreatePostHandler)
	http.HandleFunc("/post", handlers.PostPageHandler)
	http.HandleFunc("/like", handlers.LikeHandler)
	http.HandleFunc("/dislike", handlers.DislikeHandler)
	http.HandleFunc("/create_comment", handlers.CreateCommentHandler)
	http.HandleFunc("/like_comment", handlers.LikeCommentHandler)
	http.HandleFunc("/dislike_comment", handlers.DislikeCommentHandler)
	http.HandleFunc("/my_posts", handlers.MyPostsHandler)
	http.HandleFunc("/liked_posts", handlers.LikedPostsHandler)
	http.HandleFunc("/auth/google/login", handlers.GoogleLoginHandler)
	http.HandleFunc("/auth/google/callback", handlers.GoogleCallbackHandler)
	http.HandleFunc("/auth/github/login", handlers.GitHubLoginHandler)
	http.HandleFunc("/auth/github/callback", handlers.GitHubCallbackHandler)
	http.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui"))))
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	log.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
