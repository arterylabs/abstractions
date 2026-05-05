package db

import (
	"abstractions/utils/context_utils"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func WrapInTransaction(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			// if the request is a GET or options request, we don't need to wrap it in a transaction
			if c.Request().Method == "GET" || c.Request().Method == "OPTIONS" {
				c = context_utils.SetTransactionInContext(c, db)
				return next(c)
			}

			tx := db.Begin()
			if tx.Error != nil {
				return tx.Error
			}
			c = context_utils.SetTransactionInContext(c, tx)
			err := next(c)
			if err != nil {
				tx.Rollback()
				return err
			}
			return tx.Commit().Error
		}
	}
}
