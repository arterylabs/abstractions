package seeders

import (
	"abstractions/logger"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Upsert inserts or updates a single record using the provided conflict (key) columns.
// value should be a pointer to a struct (GORM model). updateColumns are the columns that will be updated on conflict.
func Upsert(ctx context.Context, db *gorm.DB, value any, conflictColumns []string, updateColumns []string) error {
	if len(conflictColumns) == 0 {
		return nil // nothing to do; silently ignore (caller responsibility to pass key)
	}
	cols := make([]clause.Column, 0, len(conflictColumns))
	for _, c := range conflictColumns {
		cols = append(cols, clause.Column{Name: c})
	}
	res := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(value)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// UpsertBatch performs an upsert for multiple records of the same model type.
// values must be a slice of pointers to structs. All share same conflict & update columns.
func UpsertBatch(ctx context.Context, db *gorm.DB, values any, conflictColumns []string, updateColumns []string) error {
	if len(conflictColumns) == 0 {
		return nil
	}
	cols := make([]clause.Column, 0, len(conflictColumns))
	for _, c := range conflictColumns {
		cols = append(cols, clause.Column{Name: c})
	}
	res := db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   cols,
		DoUpdates: clause.AssignmentColumns(updateColumns),
	}).Create(values)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

// RunWithTransaction allows a seeder to run inside its own transaction (optional helper)
// If fn returns error the transaction is rolled back; otherwise committed.
func RunWithTransaction(ctx context.Context, db *gorm.DB, fn func(tx *gorm.DB) error) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error { return fn(tx) })
}

// Example of how services register seeders (to be placed in service packages):
//
//	func init() {
//		seeders.Register(&someSeeder{})
//	}
//
// type someSeeder struct{}
// func (s *someSeeder) Name() string { return "some-seeder" }
//
//	func (s *someSeeder) Seed(ctx context.Context, db *gorm.DB) error {
//		model := &models.YourModel{Key: "value", Field: "data"}
//		return Upsert(ctx, db, model, []string{"key"}, []string{"field"})
//	}
//
// IMPORTANT: Do not put placeholder/example registrations in production code; above is documentation only.
func init() {
	logger.Debug(context.Background(), "seeders helpers loaded")
}
