package context_utils

import (
	"abstractions/db/models"
	"fmt"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func GetTransactionFromContext(c echo.Context) *gorm.DB {
	tx, ok := c.Get("tx").(*gorm.DB)
	if !ok {
		return nil
	}
	return tx
}

func SetTransactionInContext(c echo.Context, tx *gorm.DB) echo.Context {
	c.Set("tx", tx)
	return c
}

// Services context utilities - using interface{} to avoid import cycle
func GetServicesFromContext(c echo.Context) interface{} {
	return c.Get("services")
}

func SetServicesInContext(c echo.Context, svcs interface{}) echo.Context {
	c.Set("services", svcs)
	return c
}

func GetUserFromContext(c echo.Context) (*models.User, error) {
	user := c.Get("user")
	if user == nil {
		err := fmt.Errorf("could not retrieve user from context")
		return nil, err
	}

	usr, ok := user.(models.User)
	if !ok {
		err := fmt.Errorf("echo.Context has wrong type")
		return nil, err
	}
	return &usr, nil
}

func GetSessionFromContext(c echo.Context) (*models.Session, error) {
	session := c.Get("session")
	if session == nil {
		err := fmt.Errorf("could not retrieve session from context")
		return nil, err
	}

	sess, ok := session.(models.Session)
	if !ok {
		err := fmt.Errorf("echo.Context has wrong session type")
		return nil, err
	}
	return &sess, nil
}

func SetUserInContext(c echo.Context, user models.User) echo.Context {
	c.Set("user", user)
	return c
}

func SetSessionInContext(c echo.Context, session models.Session) echo.Context {
	c.Set("session", session)
	return c
}

// IsSessionImpersonated checks if the current session is an impersonated session
func IsSessionImpersonated(c echo.Context) bool {
	session, err := GetSessionFromContext(c)
	if err != nil {
		return false
	}
	return session.ImpersonatedBy != nil
}
