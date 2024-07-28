package controller

import (
	"BackendVoting/types"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth/gothic"
	"log"
	"net/http"
)

func CheckUserAndInsertIntoDB(db *sql.DB, username string, provider string, accessToken string, expiresAt string) error {
	row, err := db.Query("SELECT * FROM USERVIEW WHERE USERNAME = $1 AND PROVIDER = $2 AND ACCESS_TOKEN = $3", username, provider, accessToken)
	if err != nil {
		return fmt.Errorf("error querying user: %w", err)
	}
	if row.Next() {
		return nil
	}
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}
	_, err = tx.Exec("INSERT INTO USERS (USERNAME) VALUES ($1)", username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting user: %w", err)
	}

	row, err = tx.Query("SELECT ID FROM USERS WHERE USERNAME = $1", username)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error querying user: %w", err)
	}
	var id int
	if row.Next() {
		err = row.Scan(&id)
	} else {
		tx.Rollback()
		return fmt.Errorf("error scanning user: %w", err)

	}
	_, err = tx.Exec("INSERT INTO OAUTH_CREDENTIALS (USER_ID, PROVIDER, PROVIDER_ID, ACCESS_TOKEN, EXPIRES_AT) VALUES ($1, $2, $3, $4, $5)", id, provider, username, accessToken, expiresAt)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting oauth credentials: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	return nil
}

func OauthCallback(c *gin.Context, frontendURL string) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}
	fmt.Println("Provider callback:", provider)
	user, err := gothic.CompleteUserAuth(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.SetCookie("user", user.Name, 3600, "/", frontendURL, false, true)
	c.SetCookie("AccessToken", user.AccessToken, 3600, "/", frontendURL, false, true)
	c.JSON(http.StatusOK, user)
	//c.Redirect(http.StatusFound, frontendURL)

}

func OauthBeginn(c *gin.Context) {
	provider := c.Param("provider")
	q := c.Request.URL.Query()
	q.Add("provider", provider)
	c.Request.URL.RawQuery = q.Encode()
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provider is required"})
		return
	}
	// fmt.Println("Provider:", provider)
	gothic.BeginAuthHandler(c.Writer, c.Request)
}

func PingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func getPosts(db *sql.DB) ([]types.Post, error) {
	var posts []types.Post
	rows, err := db.Query("SELECT * FROM \"POSTVIEW\"")
	if err != nil {
		return nil, fmt.Errorf("error querying posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var post types.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ShortContent, &post.Owner, &post.Timestamp, &post.OwnerUsername, &post.TotalVotes, &post.Upvotes, &post.Downvotes)
		if err != nil {
			return nil, fmt.Errorf("error scanning post: %w", err)
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return posts, nil
}

func GetPostsHandler(c *gin.Context, db *sql.DB) {
	posts, err := getPosts(db)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"posts": posts})
}

func AuthCheckHandler(c *gin.Context, db *sql.DB) {
	username := c.Param("username")
	if username == "" {
		c.JSON(400, gin.H{"error": "username is required"})
		return
	}
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(400, gin.H{"error": "provider is required"})
		return
	}
	accessToken := c.Param("accessToken")
	if accessToken == "" {
		c.JSON(400, gin.H{"error": "accessToken is required"})
		return
	}

	row, err := db.Query("SELECT * FROM USERVIEW WHERE USERNAME = $1 AND PROVIDER = $2 AND ACCESS_TOKEN = $3", username, provider, accessToken)
	if err != nil {
		log.Fatal(err)
	}
	if row.Next() {
		var id int
		var username string
		var roleId int
		var role string
		var provider string
		var accessToken string
		var expiresAt string

		err := row.Scan(&id, &username, &roleId, &role, &provider, &accessToken, &expiresAt)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(200, gin.H{
			"message":   "success",
			"username":  username,
			"provider":  provider,
			"role":      role,
			"expiresAt": expiresAt,
		})
		return
	} else {
		c.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
}
