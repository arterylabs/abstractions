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

	redisOpt := asynq.RedisClientOpt{Addr: config.Environment.RedisHost + ":" + strconv.Itoa(config.Environment.RedisPort)}

	enqueuer := asynq.NewClient(redisOpt)
	defer enqueuer.Close()

	srv := asynq.NewServer(
		redisOpt,
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
	jobs.AddHandlers(mux, db, config, enqueuer)

	// Periodic scheduler: enqueues newsletter_scheduler every hour at minute 0.
	scheduler := asynq.NewScheduler(redisOpt, nil)
	newsletterTask, err := jobs.NewNewsletterSchedulerTask()
	if err != nil {
		log.Fatalf("failed to build newsletter scheduler task: %v", err)
	}
	if _, err := scheduler.Register("0 * * * *", newsletterTask); err != nil {
		log.Fatalf("failed to register newsletter cron: %v", err)
	}
	if err := scheduler.Start(); err != nil {
		log.Fatalf("could not start scheduler: %v", err)
	}
	defer scheduler.Shutdown()

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
