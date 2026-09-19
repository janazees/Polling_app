package handlers

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"polling-app/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserCollection *mongo.Collection
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c *gin.Context) {

	var request SignupRequest

	// Read the JSON sent by the frontend.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Clean the email before validation.
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	// Validate required fields.
	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	// Require a basic minimum password length.
	if len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Password must be at least 8 characters",
		})
		return
	}

	// Check whether the email is already registered.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var existingUser models.User

	err := h.UserCollection.FindOne(
		ctx,
		bson.M{"email": request.Email},
	).Decode(&existingUser)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email is already registered",
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not check existing user",
		})
		return
	}

	// Hash the password before storing it.
	passwordHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not secure password",
		})
		return
	}

	// Create the user document.
	user := models.User{
		ID:           bson.NewObjectID(),
		Email:        request.Email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	// Save the user to MongoDB.
	_, err = h.UserCollection.InsertOne(ctx, user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create account",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
	})
}

func (h *AuthHandler) Login(c *gin.Context) {

	var request LoginRequest

	// Read JSON sent by the client.
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	// Clean the email.
	request.Email = strings.TrimSpace(strings.ToLower(request.Email))

	// Validate required fields.
	if request.Email == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email and password are required",
		})
		return
	}

	// Create a timeout for the database operation.
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Find the user by email.
	var user models.User

	err := h.UserCollection.FindOne(
		ctx,
		bson.M{"email": request.Email},
	).Decode(&user)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not find user",
		})
		return
	}

	// Compare the entered password with the stored bcrypt hash.
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(request.Password),
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Get JWT secret from .env
	jwtSecret := os.Getenv("JWT_SECRET")

	if jwtSecret == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "JWT secret is not configured",
		})
		return
	}

	// Create JWT claims.
	claims := jwt.MapClaims{
		"userId": user.ID.Hex(),
		"email":  user.Email,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	// Create the JWT.
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Sign the JWT.
	signedToken, err := token.SignedString([]byte(jwtSecret))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create login token",
		})
		return
	}

	// Return the token to the client.
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   signedToken,
	})
}