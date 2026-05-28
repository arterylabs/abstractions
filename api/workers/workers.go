package main

import (
	app_notifications "abstractions/clients/app_notifications"
	"abstractions/config"
	"abstractions/db"
	"abstractions/logger"
	"abstractions/workers/jobs"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hibiken/asynq"
)

func logStartupHealth(cfg *config.Config, workerConfig asynq.Config) {
	ctx := context.Background()
	env := cfg.Environment

	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	logger.Info(ctx, "🔄 Workers starting up")
	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	logger.Info(ctx, fmt.Sprintf("🌍  Environment : %s", env.ENV))
	logger.Info(ctx, fmt.Sprintf("🗄️   Database    : %s@%s:%d/%s (sslmode=%s)",
		env.DatabaseUser, env.DatabaseHost, env.DatabasePort, env.DatabaseName, env.DatabaseSSLMode))
	logger.Info(ctx, fmt.Sprintf("🔴  Redis       : %s:%d", env.RedisHost, env.RedisPort))

	queues := ""
	for q, p := range workerConfig.Queues {
		if queues != "" {
			queues += " "
		}
		queues += fmt.Sprintf("%s(%d)", q, p)
	}
	logger.Info(ctx, fmt.Sprintf("📋  Queues      : %s", queues))
	logger.Info(ctx, fmt.Sprintf("🔧  Concurrency : %d", workerConfig.Concurrency))

	logger.Info(ctx, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

func main() {
	ctx := context.Background()

	logger.Info(ctx, "loading config")
	cfg := config.Load()
	logger.Info(ctx, "config loaded ", cfg.Environment.DatabaseName)

	dbConn := db.Connect(ctx, cfg)

	redisOpt := asynq.RedisClientOpt{Addr: cfg.Environment.RedisHost + ":" + strconv.Itoa(cfg.Environment.RedisPort)}

	enqueuer := jobs.NewClient(cfg.Environment.RedisHost + ":" + strconv.Itoa(cfg.Environment.RedisPort))
	defer enqueuer.Close()

	workerConfig := asynq.Config{
		Concurrency: 10,
		Queues: map[string]int{
			jobs.QueueCritical: 6,
			jobs.QueueDefault:  3,
			jobs.QueueLow:      1,
		},
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			if asynq.IsPanicError(err) {
				go app_notifications.SendPanicNotification(context.Background(), app_notifications.PanicParams{
					Source:  fmt.Sprintf("job:%s", task.Type()),
					Message: err.Error(),
				})
			} else {
				go app_notifications.SendJobFailureNotification(context.Background(), app_notifications.ErrorParams{
					Source:  fmt.Sprintf("job:%s", task.Type()),
					Message: err.Error(),
				})
			}
		}),
	}

	logStartupHealth(cfg, workerConfig)

	srv := asynq.NewServer(redisOpt, workerConfig)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	jobs.AddHandlers(mux, dbConn, cfg, enqueuer)

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
