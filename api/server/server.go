package server

import (
	"abstractions/clients/email"
	"abstractions/config"
	"abstractions/db"
	"abstractions/db/seeders"
	"context"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	Echo   *echo.Echo
	DB     *gorm.DB
	Config *config.Config
}

func NewServer(cfg *config.Config) *Server {
	// establish db connection first
	dbConn := db.Connect(context.Background(), cfg)

	// Register the email suppression store on the root connection. Emails are sent from
	// fire-and-forget goroutines that outlive the request transaction, so the store must
	// use this root connection rather than a per-request transaction.
	email.SetSuppressionStore(email.NewGormSuppressionStore(dbConn))

	// run global seeders (safe / idempotent)
	seeders.Run(context.Background(), dbConn)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	return &Server{
		Echo:   e,
		DB:     dbConn,
		Config: cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(":" + addr)
}
