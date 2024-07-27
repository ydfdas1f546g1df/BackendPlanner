package main

import (
	"BackendVoting/dbt"
	"BackendVoting/types"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var db *sql.DB
var dbUser = "myuser"
var dbPassword = "mypass"
var dbName = "votingksh"
var dbHost = "localhost"
var frontendURL = "http://localhost:5173"

func getPosts(db *sql.DB) ([]types.Post, error) {
	var posts []types.Post
	rows, err := db.Query("SELECT * FROM \"POSTVIEW\"")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post types.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortContent, &post.Owner, &post.Timestamp, &post.OwnerUsername, &post.TotalVotes, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
func getPostsHandler(c *gin.Context) {
	posts, err := getPosts(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"posts": posts})
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
	db = dbt.InitDB(dbUser, dbPassword, dbName, dbHost)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/posts", getPostsHandler)

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		return
	}
}
