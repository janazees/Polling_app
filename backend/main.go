package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"polling-app/backend/database"
	"polling-app/backend/handlers"
	"polling-app/backend/middleware"
	"polling-app/backend/realtime"
	"polling-app/backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// allowedOrigins returns the list of frontend origins allowed to call this API.
// Set ALLOWED_ORIGINS (comma-separated) in your hosting environment to override.
func allowedOrigins() []string {
	defaults := []string{
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"https://polling-app-amber-omega.vercel.app",
	}

	extra := os.Getenv("ALLOWED_ORIGINS")
	if extra == "" {
		return defaults
	}

	for _, origin := range strings.Split(extra, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			defaults = append(defaults, origin)
		}
	}
	return defaults
}

func main() {

	// Connect to MongoDB when the server starts.
	client, err := database.ConnectMongoDB()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("MongoDB connected successfully")

	// Connect to Redis when the server starts.
	redisClient, err := database.ConnectRedis()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Redis connected successfully")

	// Select the database used by the polling application.
	db := client.Database("polling_app")

	// Get the collections.
	userCollection := db.Collection("users")
	pollCollection := db.Collection("polls")
	voteCollection := db.Collection("votes")

	realtimeManager := realtime.NewManager()

	// Create the authentication handler.
	authHandler := &handlers.AuthHandler{
		UserCollection: userCollection,
	}

	// Create the poll handler.
	pollHandler := &handlers.PollHandler{
		PollCollection:  pollCollection,
		VoteCollection:  voteCollection,
		RedisClient:     redisClient,
		RealtimeManager: realtimeManager,
	}

	// Create Gin router.
	r := gin.Default()

	origins := allowedOrigins()

	// Allow the React frontend to communicate with the Go backend.
	r.Use(cors.New(cors.Config{
		AllowOrigins: origins,
		// Also allow Vercel preview deployments (e.g. polling-app-git-main-xyz.vercel.app).
		AllowOriginFunc: func(origin string) bool {
			return strings.HasSuffix(origin, ".vercel.app") && strings.HasPrefix(origin, "https://polling-app")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health-check endpoint.
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Polling backend is running",
		})
	})

	// Register authentication routes.
	routes.SetupAuthRoutes(r, authHandler)

	// Register poll routes.
	routes.SetupPollRoutes(r, pollHandler)

	// Protected test route.
	protected := r.Group("/api/protected")
	protected.Use(middleware.AuthMiddleware())

	protected.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "You are authenticated!",
		})
	})

	// Hosting platforms (Render, Railway, Fly) provide PORT.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Backend running on port " + port)
	log.Println("Allowed origins:", origins)

	// Start the API server.
	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
