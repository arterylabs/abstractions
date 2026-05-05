package main

import (
	application "abstractions"
	"abstractions/config"
	"abstractions/db"
	"os"
	"strconv"
)

func main() {
	cfg := config.Load()

	// if the user runs main.go with the migrate argument and steps we run the migrations
	// example: go run cmd/main.go migrate
	args := os.Args[1:]
	if len(args) > 0 {
		if args[0] == "migrate" {
			// Create a copy of the config for migration
			migrationEnv := *cfg.Environment

			// Check for database configuration flags
			steps := 0
			i := 1
			for i < len(args) {
				switch args[i] {
				case "--host":
					if i+1 < len(args) {
						migrationEnv.DatabaseHost = args[i+1]
						i += 2
					}
				case "--port":
					if i+1 < len(args) {
						if port, err := strconv.Atoi(args[i+1]); err == nil {
							migrationEnv.DatabasePort = port
						}
						i += 2
					}
				case "--user":
					if i+1 < len(args) {
						migrationEnv.DatabaseUser = args[i+1]
						i += 2
					}
				case "--password":
					if i+1 < len(args) {
						migrationEnv.DatabasePassword = args[i+1]
						i += 2
					}
				case "--dbname":
					if i+1 < len(args) {
						migrationEnv.DatabaseName = args[i+1]
						i += 2
					}
				case "--sslmode":
					if i+1 < len(args) {
						migrationEnv.DatabaseSSLMode = args[i+1]
						i += 2
					}
				case "force":
					versionToForce := 0
					if i+1 < len(args) {
						versionToForce, _ = strconv.Atoi(args[i+1])
					}
					db.ForceMigrationVersion(versionToForce)
					return
				case "up":
					db.MigrateWithConfig(&migrationEnv, 1)
					return
				case "down":
					db.MigrateWithConfig(&migrationEnv, -1)
					return
				default:
					if val, err := strconv.Atoi(args[i]); err == nil {
						steps = val
					}
					i++
				}
			}

			// Run migration with specified or default steps
			db.MigrateWithConfig(&migrationEnv, steps)
			return
		}
	}

	application.Start(cfg)
}
