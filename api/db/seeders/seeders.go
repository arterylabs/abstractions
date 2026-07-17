package seeders

import (
	"abstractions/logger"
	"context"

	"gorm.io/gorm"
)

// Seeder defines the contract each seeder should implement.
// Key is used for idempotency (upsert behavior) and should uniquely identify the seeded record.
type Seeder interface {
	// Name returns a short name of the seeder (used in logs)
	Name() string
	// Seed executes the seeding logic. It should be idempotent. If the record already exists, it should update it.
	Seed(ctx context.Context, db *gorm.DB) error
}

// Registry maintains a list of seeders to run.
type Registry struct {
	seeders []Seeder
}

// NewRegistry creates a new seeder registry.
func NewRegistry() *Registry {
	return &Registry{seeders: []Seeder{}}
}

// Register appends seeders to the registry (typically from service packages during init()).
func (r *Registry) Register(seeders ...Seeder) {
	for _, s := range seeders {
		if s == nil {
			continue
		}
		r.seeders = append(r.seeders, s)
	}
}

// Run executes all registered seeders in order.
func (r *Registry) Run(ctx context.Context, db *gorm.DB) {
	for _, s := range r.seeders {
		logger.Info(ctx, "running seeder", s.Name())
		if err := s.Seed(ctx, db); err != nil {
			logger.Error(ctx, "seeder failed", s.Name(), err)
		} else {
			logger.Info(ctx, "seeder completed", s.Name())
		}
	}
}

// defaultRegistry is the global registry used by the application.
var defaultRegistry = NewRegistry()

// Register registers seeders in the global registry.
func Register(seeders ...Seeder) { defaultRegistry.Register(seeders...) }

// Run runs all seeders in the global registry.
func Run(ctx context.Context, db *gorm.DB) { defaultRegistry.Run(ctx, db) }
