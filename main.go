package main

import (
	"BackendVoting/dbt"
	"BackendVoting/types"
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var dbUser = "myuser"
var dbPassword = "mypass"
var dbName = "votingksh"
var dbHost = "localhost"
var frontendURL = "http://localhost:5173"

func getPosts(db *sql.DB) gin.H {
	var posts []types.Post
	rows := dbt.QueryDB(db, "SELECT * FROM posts")
	for rows.Next() {
		var post types.Post
		rows.Scan(&post.ID, &post.Title, &post.Content)
		posts = append(posts, post)
	}
	return posts
}

func main() {
	r := gin.Default()
	config := cors.Config{
		AllowOrigins:     []string{frontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))
	db := dbt.InitDB(dbUser, dbPassword, dbName, dbHost)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/posts", func(c *gin.Context) {
		c.JSON(200, GetPosts())
	})

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
