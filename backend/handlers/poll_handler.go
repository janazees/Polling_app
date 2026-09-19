package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"polling-app/backend/models"
	"polling-app/backend/realtime"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type PollHandler struct {
	PollCollection  *mongo.Collection
	VoteCollection  *mongo.Collection
	RedisClient     *redis.Client
	RealtimeManager *realtime.Manager
}

type CreatePollRequest struct {
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

type VoteRequest struct {
	OptionID string `json:"optionId"`
	VoterID  string `json:"voterId"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Development only.
		return true
	},
}

// --------------------------------------------------
// CREATE POLL
// --------------------------------------------------

func (h *PollHandler) CreatePoll(c *gin.Context) {

	var request CreatePollRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data",
		})
		return
	}

	request.Question = strings.TrimSpace(request.Question)

	if request.Question == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Question is required",
		})
		return
	}

	if len(request.Question) > 200 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Question must be 200 characters or less",
		})
		return
	}

	if len(request.Options) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "A poll must have at least 2 options",
		})
		return
	}

	if len(request.Options) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "A poll can have at most 10 options",
		})
		return
	}

	options := make([]models.PollOption, 0, len(request.Options))

	for _, optionText := range request.Options {

		optionText = strings.TrimSpace(optionText)

		if optionText == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Options cannot be empty",
			})
			return
		}

		if len(optionText) > 100 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Each option must be 100 characters or less",
			})
			return
		}

		options = append(options, models.PollOption{
			ID:   uuid.NewString(),
			Text: optionText,
		})
	}

	tokenValue, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User authentication required",
		})
		return
	}

	token, ok := tokenValue.(*jwt.Token)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		return
	}

	userIDString, ok := claims["userId"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID missing from token",
		})
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	shareCode := uuid.NewString()

	poll := models.Poll{
		ID:        bson.NewObjectID(),
		Question:  request.Question,
		Options:   options,
		CreatedBy: userID,
		ShareCode: shareCode,
		IsClosed:  false,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	_, err = h.PollCollection.InsertOne(ctx, poll)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not create poll",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Poll created successfully",
		"poll":    poll,
	})
}

// --------------------------------------------------
// GET MY POLLS
// --------------------------------------------------

func (h *PollHandler) GetMyPolls(c *gin.Context) {

	tokenValue, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User authentication required",
		})
		return
	}

	token, ok := tokenValue.(*jwt.Token)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		return
	}

	userIDString, ok := claims["userId"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID missing from token",
		})
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	cursor, err := h.PollCollection.Find(
		ctx,
		bson.M{
			"createdBy": userID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch your polls",
		})
		return
	}

	defer cursor.Close(ctx)

	var polls []models.Poll

	if err := cursor.All(ctx, &polls); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read your polls",
		})
		return
	}

	if polls == nil {
		polls = []models.Poll{}
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
	})
}

// --------------------------------------------------
// DELETE POLL
// --------------------------------------------------

func (h *PollHandler) DeletePoll(c *gin.Context) {

	shareCode := c.Param("shareCode")

	if shareCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Poll share code is required",
		})
		return
	}

	tokenValue, exists := c.Get("user")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User authentication required",
		})
		return
	}

	token, ok := tokenValue.(*jwt.Token)

	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication token",
		})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		return
	}

	userIDString, ok := claims["userId"].(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID missing from token",
		})
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDString)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	result, err := h.PollCollection.DeleteOne(
		ctx,
		bson.M{
			"shareCode": shareCode,
			"createdBy": userID,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete poll",
		})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Poll not found or you do not have permission to delete it",
		})
		return
	}

	if h.RedisClient != nil {

		redisKey := "poll:" + shareCode + ":votes"

		if err := h.RedisClient.Del(
			context.Background(),
			redisKey,
		).Err(); err != nil {
			log.Println("Failed to delete Redis vote data:", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll deleted successfully",
	})
}

// --------------------------------------------------
// GET POLL
// --------------------------------------------------

func (h *PollHandler) GetPollByShareCode(c *gin.Context) {

	shareCode := strings.TrimSpace(c.Param("shareCode"))

	if shareCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Share code is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var poll models.Poll

	err := h.PollCollection.FindOne(
		ctx,
		bson.M{"shareCode": shareCode},
	).Decode(&poll)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Poll not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch poll",
		})
		return
	}

	// ----------------------------------------------
	// GET VOTE COUNTS
	// ----------------------------------------------

	results := make(map[string]int)

	for _, option := range poll.Options {
		results[option.ID] = 0
	}

	if h.RedisClient != nil {

		redisKey := "poll:" + shareCode + ":votes"

		voteCounts, err := h.RedisClient.HGetAll(
			ctx,
			redisKey,
		).Result()

		if err == nil && len(voteCounts) > 0 {

			for optionID, value := range voteCounts {

				var count int

				_, err := fmt.Sscanf(value, "%d", &count)

				if err == nil {
					results[optionID] = count
				}
			}

		} else {

			cursor, err := h.VoteCollection.Find(
				ctx,
				bson.M{"pollId": poll.ID},
			)

			if err == nil {

				defer cursor.Close(ctx)

				for cursor.Next(ctx) {

					var vote models.Vote

					if err := cursor.Decode(&vote); err != nil {
						continue
					}

					results[vote.OptionID]++
				}

				for optionID, count := range results {

					if count > 0 {
						_ = h.RedisClient.HSet(
							ctx,
							redisKey,
							optionID,
							count,
						).Err()
					}
				}
			}
		}
	}

	totalVotes := 0

	for _, count := range results {
		totalVotes += count
	}

	c.JSON(http.StatusOK, gin.H{
		"poll":       poll,
		"results":    results,
		"totalVotes": totalVotes,
	})
}

// --------------------------------------------------
// VOTE
// --------------------------------------------------

func (h *PollHandler) Vote(c *gin.Context) {

	shareCode := strings.TrimSpace(c.Param("shareCode"))

	if shareCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Share code is required",
		})
		return
	}

	var request VoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	request.OptionID = strings.TrimSpace(request.OptionID)

	if request.OptionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Option ID is required",
		})
		return
	}

	request.VoterID = strings.TrimSpace(request.VoterID)

	if request.VoterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Voter ID is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var poll models.Poll

	err := h.PollCollection.FindOne(
		ctx,
		bson.M{"shareCode": shareCode},
	).Decode(&poll)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Poll not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to find poll",
		})
		return
	}

	// ----------------------------------------------
	// CHECK OPTION
	// ----------------------------------------------

	optionExists := false

	for _, option := range poll.Options {

		if option.ID == request.OptionID {
			optionExists = true
			break
		}
	}

	if !optionExists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid option for this poll",
		})
		return
	}

	// ----------------------------------------------
	// CHECK IF VOTER ALREADY VOTED
	// ----------------------------------------------

	var existingVote models.Vote

	err = h.VoteCollection.FindOne(
		ctx,
		bson.M{
			"pollId":  poll.ID,
			"voterId": request.VoterID,
		},
	).Decode(&existingVote)

	if err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "You have already voted in this poll",
		})
		return
	}

	if err != mongo.ErrNoDocuments {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check previous vote",
		})
		return
	}

	// ----------------------------------------------
	// SAVE VOTE TO MONGODB
	// ----------------------------------------------

	vote := models.Vote{
		PollID:    poll.ID,
		OptionID:  request.OptionID,
		VoterID:   request.VoterID,
		CreatedAt: time.Now(),
	}

	_, err = h.VoteCollection.InsertOne(ctx, vote)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save vote",
		})
		return
	}

	// ----------------------------------------------
	// UPDATE REDIS COUNT
	// ----------------------------------------------

	redisKey := "poll:" + shareCode + ":votes"

	newCount, err := h.RedisClient.HIncrBy(
		ctx,
		redisKey,
		request.OptionID,
		1,
	).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update live vote count",
		})
		return
	}

	// ----------------------------------------------
	// GET ALL CURRENT COUNTS
	// ----------------------------------------------

	voteCounts, err := h.RedisClient.HGetAll(
		ctx,
		redisKey,
	).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get vote counts",
		})
		return
	}

	results := make(map[string]int)

	for _, option := range poll.Options {
		results[option.ID] = 0
	}

	for optionID, value := range voteCounts {

		var count int

		_, err := fmt.Sscanf(value, "%d", &count)

		if err == nil {
			results[optionID] = count
		}
	}

	// ----------------------------------------------
	// CALCULATE TOTAL
	// ----------------------------------------------

	totalVotes := 0

	for _, count := range results {
		totalVotes += count
	}

	// ----------------------------------------------
	// PUBLISH LIVE UPDATE
	// ----------------------------------------------

	message := map[string]interface{}{
		"pollId":     poll.ID.Hex(),
		"results":    results,
		"totalVotes": totalVotes,
	}

	messageJSON, err := json.Marshal(message)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create live update",
		})
		return
	}

	channel := "poll:" + shareCode + ":updates"

	err = h.RedisClient.Publish(
		ctx,
		channel,
		messageJSON,
	).Err()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to publish live update",
		})
		return
	}

	// ----------------------------------------------
	// RESPONSE
	// ----------------------------------------------

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Vote recorded successfully",
		"count":      newCount,
		"results":    results,
		"totalVotes": totalVotes,
	})
}

// --------------------------------------------------
// WEBSOCKET
// --------------------------------------------------

func (h *PollHandler) PollWebSocket(c *gin.Context) {

	shareCode := strings.TrimSpace(c.Param("shareCode"))

	if shareCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Share code is required",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	var poll models.Poll

	err := h.PollCollection.FindOne(
		ctx,
		bson.M{"shareCode": shareCode},
	).Decode(&poll)

	if err != nil {

		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Poll not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to find poll",
		})
		return
	}

	conn, err := upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)

	if err != nil {
		return
	}

	client := &realtime.Client{
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	hub := h.RealtimeManager.GetHub(shareCode)

	hub.AddClient(client)

	h.RealtimeManager.StartRedisSubscriber(
		h.RedisClient,
		shareCode,
	)

	defer func() {
		hub.RemoveClient(client)
		conn.Close()
	}()

	go func() {

		for message := range client.Send {

			err := conn.WriteMessage(
				websocket.TextMessage,
				message,
			)

			if err != nil {
				return
			}
		}
	}()

	for {

		_, _, err := conn.ReadMessage()

		if err != nil {
			break
		}
	}
}

// --------------------------------------------------
// CLOSING POLL
// --------------------------------------------------

func (h *PollHandler) ClosePoll(c *gin.Context) {
	shareCode := c.Param("shareCode")

	// Get the JWT token stored by the existing middleware.
	userValue, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	token, ok := userValue.(*jwt.Token)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid authentication token",
		})
		return
	}

	// Read the userId from the JWT claims.
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		return
	}

	userIDString, ok := claims["userId"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in token",
		})
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Find the poll and make sure it belongs to this user.
	var poll models.Poll

	err = h.PollCollection.FindOne(
		ctx,
		bson.M{
			"shareCode": shareCode,
			"createdBy": userID,
		},
	).Decode(&poll)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Poll not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to find poll",
		})
		return
	}

	// Don't close an already closed poll.
	if poll.IsClosed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Poll is already closed",
		})
		return
	}

	// Mark the poll as closed.
	_, err = h.PollCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":       poll.ID,
			"createdBy": userID,
		},
		bson.M{
			"$set": bson.M{
				"isClosed": true,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to close poll",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Poll closed successfully",
	})
}
