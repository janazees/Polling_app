package routes

import (
	"polling-app/backend/handlers"
	"polling-app/backend/middleware"

	"github.com/gin-gonic/gin"
)

func SetupPollRoutes(
	r *gin.Engine,
	pollHandler *handlers.PollHandler,
) {

	// -----------------------------------------
	// PUBLIC POLL ROUTES
	// -----------------------------------------

	r.GET(
		"/api/polls/:shareCode",
		pollHandler.GetPollByShareCode,
	)

	r.POST(
		"/api/polls/:shareCode/vote",
		pollHandler.Vote,
	)

	r.GET(
		"/api/polls/:shareCode/ws",
		pollHandler.PollWebSocket,
	)

	// -----------------------------------------
	// PROTECTED POLL ROUTES
	// -----------------------------------------

	pollRoutes := r.Group("/api/polls")
	pollRoutes.Use(middleware.AuthMiddleware())

	// Create a poll
	pollRoutes.POST(
		"",
		pollHandler.CreatePoll,
	)

	// Get polls created by logged-in user
	pollRoutes.GET(
		"",
		pollHandler.GetMyPolls,
	)

	// Delete a poll
	pollRoutes.DELETE(
		"/:shareCode",
		pollHandler.DeletePoll,
	)

	pollRoutes.POST(
		"/:shareCode/close",
		pollHandler.ClosePoll,
	)
}
