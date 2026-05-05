package context_utils

import (
	"abstractions/db/models"
	"fmt"

	"github.com/google/uuid"
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

// GetImpersonatorID returns the ID of the user who is impersonating, if any
func GetImpersonatorID(c echo.Context) *uuid.UUID {
	session, err := GetSessionFromContext(c)
	if err != nil {
		return nil
	}
	return session.ImpersonatedBy
}

// func GetTransactionFromContext(ctx context.Context) (*pop.Connection, error) {
// 	tx := ctx.Value("tx")
// 	if tx == nil {
// 		return nil, errors.New("transaction not found")
// 	}

// 	popTx, ok := tx.(*pop.Connection)
// 	if !ok {
// 		return nil, errors.New("transaction malformed in context")
// 	}

// 	return popTx, nil
// 	// echoContext, err := GetEchoContextFromContext(ctx)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// return GetTransactionFromEchoContext(echoContext)
// }

// func GetTransactionFromEchoContext(c echo.Context) (*pop.Connection, error) {
// 	tx := c.Get("tx")
// 	if tx == nil {
// 		return nil, errors.New("transaction not found")
// 	}

// 	popTx, ok := tx.(*pop.Connection)
// 	if !ok {
// 		return nil, errors.New("transaction malformed in context")
// 	}

// 	return popTx, nil
// }

// func GetOrganisationIDFromContext(ctx context.Context) (string, error) {
// 	orgId := ctx.Value("organisation_id")
// 	if orgId == nil {
// 		return "", errors.New("organisation_id not found")
// 	}

// 	orgIdString, ok := orgId.(string)
// 	if !ok {
// 		return "", errors.New("organisation_id malformed in context")
// 	}

// 	return orgIdString, nil
// }

// func AddDefaultFieldsToEchoContext(c echo.Context) {
// 	loggerDefaultFields := make(map[string]interface{})
// 	loggerDefaultFields["request_id"] = c.Response().Header().Get("X-Request-ID")
// 	loggerDefaultFields["real_ip"] = c.RealIP()
// 	loggerDefaultFields["user_agent"] = c.Request().UserAgent()
// 	loggerDefaultFields["host"] = c.Request().Host
// 	loggerDefaultFields["method"] = c.Request().Method
// 	path := c.Request().URL.EscapedPath()
// 	if path == "" {
// 		path = "/"
// 	}
// 	loggerDefaultFields["path"] = path
// 	loggerDefaultFields["query"] = c.Request().URL.RawQuery
// 	loggerDefaultFields["referer"] = c.Request().Referer()
// 	loggerDefaultFields["status"] = c.Response().Status

// 	// if the request has a token then decode it and add the user id to context so we log it
// 	sessionToken := ""
// 	if c.Request().Header.Get("X-Access-Token") != "" {
// 		sessionToken = c.Request().Header.Get("X-Access-Token")
// 	} else {
// 		tokenCookie, _ := c.Cookie("token")
// 		if tokenCookie != nil {
// 			sessionToken = tokenCookie.Value
// 		} else {
// 			sessionToken = ""
// 		}
// 	}
// 	if sessionToken != "" {
// 		userID, _, _ := DecodeJWT(sessionToken)
// 		loggerDefaultFields["user_id"] = userID
// 		if userID != "" {
// 			user, err := models.GetUserFromCache(userID)
// 			if err == nil {
// 				loggerDefaultFields["user_email"] = user.Email
// 				loggerDefaultFields["user_username"] = user.Username
// 				loggerDefaultFields["user_default_organisation"] = user.DefaultOrganisation
// 			}
// 		}
// 	}

// 	loggerDefaultFields["content_length"] = c.Request().ContentLength
// 	loggerDefaultFields["protocol"] = c.Request().Proto
// 	loggerDefaultFields["response_size"] = c.Response().Size
// 	loggerDefaultFields["start_time"] = time.Now().UTC()

// 	context := c.Request().Context()
// 	logger.AddFieldsToContext(context, loggerDefaultFields)
// 	c.SetRequest(c.Request().WithContext(context))
// }

// func GetDefaultFieldsFromEchoContext(c echo.Context, start_time time.Time) map[string]interface{} {
// 	loggerDefaultFields := make(map[string]interface{})
// 	loggerDefaultFields["request_id"] = c.Response().Header().Get(echo.HeaderXRequestID)
// 	loggerDefaultFields["real_ip"] = c.RealIP()
// 	loggerDefaultFields["user_agent"] = c.Request().UserAgent()
// 	loggerDefaultFields["host"] = c.Request().Host
// 	loggerDefaultFields["method"] = c.Request().Method
// 	path := c.Request().URL.EscapedPath()
// 	if path == "" {
// 		path = "/"
// 	}
// 	loggerDefaultFields["path"] = path
// 	loggerDefaultFields["query"] = c.Request().URL.RawQuery
// 	loggerDefaultFields["referer"] = c.Request().Referer()
// 	loggerDefaultFields["status"] = c.Response().Status

// 	// if the request has a token then decode it and add the user id to context so we log it
// 	sessionToken := ""
// 	if c.Request().Header.Get("X-Access-Token") != "" {
// 		sessionToken = c.Request().Header.Get("X-Access-Token")
// 	} else {
// 		tokenCookie, _ := c.Cookie("token")
// 		if tokenCookie != nil {
// 			sessionToken = tokenCookie.Value
// 		} else {
// 			sessionToken = ""
// 		}
// 	}
// 	if sessionToken != "" {
// 		userID, _, _ := DecodeJWT(sessionToken)
// 		loggerDefaultFields["user_id"] = userID
// 		if userID != "" {
// 			user, err := models.GetUserFromCache(userID)
// 			if err == nil {
// 				loggerDefaultFields["user_email"] = user.Email
// 				loggerDefaultFields["user_username"] = user.Username
// 				loggerDefaultFields["user_default_organisation"] = user.DefaultOrganisation
// 			}
// 		}
// 	}

// 	loggerDefaultFields["content_length"] = c.Request().ContentLength
// 	loggerDefaultFields["protocol"] = c.Request().Proto
// 	loggerDefaultFields["response_size"] = c.Response().Size
// 	loggerDefaultFields["start_time"] = start_time

// 	return loggerDefaultFields
// }

// func IsUserAuthenticated(ctx context.Context) bool {
// 	c, err := GetEchoContextFromContext(ctx)
// 	if err != nil {
// 		logger.Error(ctx, err)
// 		return false
// 	}
// 	return IsUserAuthenticatedEchoCtx(c)
// }

// func IsUserAuthenticatedEchoCtx(c echo.Context) bool {
// 	userID := ""
// 	sessionToken := ""
// 	apiKey := ""

// 	if c.Request().Header.Get("X-API-Key") != "" {
// 		apiKey = c.Request().Header.Get("X-API-Key")
// 		if apiKey == "" {
// 			logger.Error(c.Request().Context(), errors.New("access token is empty"))
// 			return false
// 		}

// 		apiKeyCacheKey := "api_key_" + apiKey
// 		apiKeyStringified, err := cache.Client.GetKey(apiKeyCacheKey)
// 		if err != nil || apiKeyStringified == "" {
// 			apiKeyInDB := models.APIKey{}
// 			decodedKey, err := base64.StdEncoding.DecodeString(apiKey)
// 			if err != nil {
// 				logger.Error(c.Request().Context(), err)
// 				return false
// 			}

// 			res := strings.Split(string(decodedKey), "-")

// 			decodedOrgId, _ := base64.StdEncoding.DecodeString(res[0])

// 			err = models.DB.Where("is_revoked = false and organisation_id in (?)", string(decodedOrgId)).First(&apiKeyInDB)
// 			if err != nil {
// 				logger.Error(c.Request().Context(), errors.New("no valid API key found"))
// 				return false
// 			}

// 			decryptedKey, _ := encryption_utils.DecryptString(environment.Config.EncryptionKey32Bit, apiKeyInDB.Token)
// 			if decryptedKey != apiKey {
// 				return false
// 			}
// 			if apiKeyInDB.ExpiresAt.Valid && apiKeyInDB.ExpiresAt.Time.Before(time.Now().UTC()) {
// 				logger.Error(c.Request().Context(), errors.New("api key expired"))
// 				return false
// 			}
// 			userID = apiKeyInDB.UserID.String()
// 			apiKeyInDB.Token = "" // empty the token because we don't want to store in the cache
// 			cache.Client.SetKey(apiKeyCacheKey, apiKeyInDB.String(), time.Hour*24)
// 		} else {
// 			apiKeyFromCache := models.APIKey{}
// 			err = json.Unmarshal([]byte(apiKeyStringified), &apiKeyFromCache)
// 			if err != nil {
// 				cache.Client.DeleteKey(apiKeyCacheKey)
// 			}
// 			// if the APIKey is revoked or has expired then return error
// 			if apiKeyFromCache.IsRevoked || (apiKeyFromCache.ExpiresAt.Valid && apiKeyFromCache.ExpiresAt.Time.Before(time.Now().UTC())) {
// 				logger.Error(c.Request().Context(), errors.New("api key expired"))
// 				return false
// 			}
// 			userID = apiKeyFromCache.UserID.String()
// 		}

// 	} else {
// 		if c.Request().Header.Get("X-Access-Token") != "" {
// 			sessionToken = c.Request().Header.Get("X-Access-Token")
// 		} else {
// 			tokenCookie, _ := c.Cookie("token")
// 			if tokenCookie != nil {
// 				sessionToken = tokenCookie.Value
// 			} else {
// 				sessionToken = ""
// 			}
// 		}
// 		if sessionToken != "" {

// 			sessionID := ""
// 			var err error
// 			userID, sessionID, err = DecodeJWT(sessionToken)
// 			if err != nil {
// 				logger.Error(c.Request().Context(), errors.New("unable to decode token"))
// 				return false
// 			}

// 			expiresAt := time.Now()

// 			sessionCacheKey := "session_" + sessionID

// 			expiresAtString, err := cache.Client.GetKey(sessionCacheKey)
// 			if err != nil || expiresAtString == "" {
// 				session := models.Session{}
// 				if err := models.DB.Find(&session, sessionID); err != nil {
// 					logger.Error(c.Request().Context(), errors.New("no session found"))
// 					return false
// 				}
// 				expiresAt = session.ExpiresAt
// 				cache.Client.SetKey(sessionCacheKey, session.ExpiresAt.Format(time.RFC3339), time.Hour*24)
// 			} else {
// 				expiresAt, err = time.Parse(time.RFC3339, expiresAtString)
// 				if err != nil {
// 					cache.Client.DeleteKey(sessionCacheKey)
// 				}
// 			}

// 			if expiresAt.Before(time.Now().UTC()) {
// 				cache.Client.DeleteKey(sessionCacheKey)
// 				logger.Error(c.Request().Context(), errors.New("session expired"))
// 				return false
// 			}
// 		}
// 	}

// 	if userID != "" {
// 		user, err := models.GetUserFromCache(userID)
// 		if err != nil || user == nil {
// 			return false
// 		}
// 		// user := models.User{}

// 		// if err := models.DB.Where("is_active = true and is_archived = false").Find(&user, userID); err != nil {
// 		// 	logger.Error(c.Request().Context(),errors.New("user not found"))
// 		// 	return false
// 		// }
// 		c.Set("user", *user)
// 		c.Set("Token", sessionToken)
// 		c.Set("User-Id", userID)
// 		c.Set("API-Key", apiKey)
// 		sentryHub := sentryecho.GetHubFromContext(c)
// 		sentryHub.ConfigureScope(func(scope *sentry.Scope) {
// 			scope.SetUser(sentry.User{Email: user.Email})
// 		})
// 		return true
// 	} else {
// 		return false
// 	}
// }

// func UserInOrganisation(ctx context.Context, organisation_id string) (bool, string) {

// 	c, err := GetEchoContextFromContext(ctx)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false, ""
// 	}

// 	user_id := fmt.Sprintf("%v", c.Get("User-Id"))

// 	if user_id == "" {
// 		logger.Error(c.Request().Context(), errors.New("user id not resolved"))
// 		return false, ""
// 	}

// 	user, err := GetUserFromContext(ctx)
// 	if user == nil || err != nil {
// 		logger.Error(ctx, "user not available in context")
// 		return false, ""
// 	}
// 	sentry.ConfigureScope(func(scope *sentry.Scope) {
// 		scope.SetUser(sentry.User{Email: user.Email})
// 	})
// 	if organisation_id == "" {
// 		logger.Error(c.Request().Context(), errors.New("unable to get organisation id"))
// 		return false, ""
// 	}

// 	return UserAddedToOrganisation(organisation_id, user_id)
// }

// func UserAddedToOrganisation(organisation_id string, user_id string) (bool, string) {
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	if !IsValidUUID(organisation_id) {
// 		id, err := models.GetOrganisationIDFromAlias(ctx, organisation_id)
// 		if err != nil {
// 			logger.Error(ctx, err)
// 			return false, ""
// 		}
// 		organisation_id = id
// 	}

// 	userOrg := models.UserOrganisation{}
// 	cacheKey := organisation_id + "_" + user_id + "_role"
// 	role, err := cache.Client.GetKey(cacheKey)
// 	if err != nil || role == "" || !roles.IsValidRole(role) {
// 		if err := models.DB.Where(fmt.Sprintf("user_id = '%s' and organisation_id = '%s'", user_id, organisation_id)).First(&userOrg); err != nil {
// 			logger.Error(ctx, err)
// 			return false, ""
// 		}
// 		cache.Client.SetKey(cacheKey, userOrg.Role.String, time.Hour*24)
// 		return true, userOrg.Role.String
// 	}

// 	return true, role
// }

// func CheckEtag(c echo.Context, requestEtag string, hashString string, etagType string) bool {
// 	etagHashGenerated := base64.StdEncoding.EncodeToString([]byte(hashString))
// 	logger.Info(c.Request().Context(), "Etag generated", etagHashGenerated)

// 	c.Response().Header().Set(etagType+"etag", etagHashGenerated)
// 	c.Response().Header().Set("Access-Control-Expose-Headers", etagType+"etag")
// 	if etagHashGenerated == requestEtag {
// 		c.NoContent(http.StatusNotModified)
// 		return true
// 	}
// 	return false
// }
