package handlers

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/wgarcia4190/garagesale/internal/platform/database"
	"github.com/wgarcia4190/garagesale/internal/platform/web"
)

// Check has handlers to implement service orchestration.
type Check struct {
	DB *sqlx.DB
}

// Health responds with a 200 OK if the service is healthy and ready for traffic
func (c *Check) Health(writer http.ResponseWriter, request *http.Request) error {

	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(request.Context(), c.DB); err != nil {

		// If the database is not ready we will tell the client and use a 500
		// status. Do not respond by just returning an error because further up in
		// the call stack will interpret that as an unhandled error.
		health.Status = "db not ready"
		return web.Respond(request.Context(), writer, health, http.StatusInternalServerError)
	}

	health.Status = "ok"
	return web.Respond(request.Context(), writer, health, http.StatusOK)
}
