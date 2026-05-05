package main

import (
	"abstractions/config"
	"abstractions/db"
	"abstractions/logger"
	"abstractions/workers/jobs"
	"context"
	"log"
	"strconv"

	"github.com/hibiken/asynq"
)

func main() {
	ctx := context.Background()

	logger.Info(ctx, "loading config")
	config := config.Load()
	logger.Info(ctx, "config loaded ", config.Environment.DatabaseName)

	db := db.Connect(ctx, config)

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: config.Environment.RedisHost + ":" + strconv.Itoa(config.Environment.RedisPort)},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	jobs.AddHandlers(mux, db, config)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
