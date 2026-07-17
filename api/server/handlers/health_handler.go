package handlers

import (
	"abstractions/server"
	"abstractions/server/responses"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	*BaseHandler
}

func NewHealthHandler(s *server.Server) *HealthHandler {
	return &HealthHandler{BaseHandler: NewBaseHandler(s)}
}

type healthResponse struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Redis    string `json:"redis"`
}

func (h *HealthHandler) Check(c echo.Context) error {
	ctx, cancel := context.WithTimeout(c.Request().Context(), 5*time.Second)
	defer cancel()

	resp := healthResponse{
		Status:   "ok",
		Database: "ok",
		Redis:    "ok",
	}
	httpStatus := http.StatusOK

	sqlDB, err := h.server.DB.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		resp.Status = "error"
		resp.Database = "error"
		httpStatus = http.StatusInternalServerError
	}

	redisAddr := fmt.Sprintf("%s:%d", h.server.Config.Environment.RedisHost, h.server.Config.Environment.RedisPort)
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	defer rdb.Close()
	if rdb.Ping(ctx).Err() != nil {
		resp.Status = "error"
		resp.Redis = "error"
		httpStatus = http.StatusInternalServerError
	}

	return responses.Response(c, httpStatus, resp)
}
