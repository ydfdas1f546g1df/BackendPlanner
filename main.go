package main

import (
	"BackendVoting/controller"
	"BackendVoting/dbt"
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"os"
)

var db *sql.DB
var GEdbUser string
var GEdbPassword string
var GEdbName string
var GEdbHost string
var GEfrontendURL string
var GEsessionSecret string

func setENV() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	requiredEnvVars := []string{
		"DB_USER", "DB_PASSWORD", "DB_NAME", "DB_HOST",
		"FRONTEND_URL", "SESSION_SECRET",
		"GITHUB_CLIENT_ID", "GITHUB_CLIENT_SECRET", "GITHUB_CLIENT_CALLBACK_URL",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			fmt.Printf("%s not set\n", envVar)
			os.Exit(1)
		}
	}

	GEdbUser = os.Getenv("DB_USER")
	GEdbPassword = os.Getenv("DB_PASSWORD")
	GEdbName = os.Getenv("DB_NAME")
	GEdbHost = os.Getenv("DB_HOST")
	GEfrontendURL = os.Getenv("FRONTEND_URL")
	GEsessionSecret = os.Getenv("SESSION_SECRET")

	err = os.Setenv("SESSION_SECRET", GEsessionSecret)
	if err != nil {
		fmt.Println("Error setting session secret")
		os.Exit(1)
	}
}

func gothInit() {
	goth.UseProviders(
		github.New(
			os.Getenv("GITHUB_CLIENT_ID"),
			os.Getenv("GITHUB_CLIENT_SECRET"),
			os.Getenv("GITHUB_CLIENT_CALLBACK_URL"),
		),
	)
}

func main() {
	setENV()
	gothInit()

	r := gin.Default()

	store := cookie.NewStore([]byte(GEsessionSecret))
	gothic.Store = store
	r.Use(sessions.Sessions("gothic", store))

	config := cors.Config{
		AllowOrigins:     []string{GEfrontendURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	db = dbt.InitDB(GEdbUser, GEdbPassword, GEdbName, GEdbHost)

	r.GET("/ping", controller.PingHandler)

	r.GET("/posts", func(c *gin.Context) {
		controller.GetPostsHandler(c, db)
	})

	r.POST("/authcheck", func(c *gin.Context) {
		controller.AuthCheckHandler(c, db)
	})

	r.GET("/auth/:provider", controller.OauthBeginn)

	r.GET("/auth/:provider/callback", func(c *gin.Context) {
		controller.OauthCallback(c, GEfrontendURL)
	})

	r.GET("/", func(c *gin.Context) {
		c.Redirect(302, GEfrontendURL)
	})

	err := r.Run("0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
