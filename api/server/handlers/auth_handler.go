package handlers

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	integration_service "abstractions/services/integrations"
	github_integration "abstractions/services/integrations/core/github"
	user_service "abstractions/services/user"
	"abstractions/utils/context_utils"
	"abstractions/utils/encryption_utils"
	"abstractions/utils/jwt_utils"
	"abstractions/utils/password_utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	*BaseHandler
}

// oauthHTTPClient is a package-level HTTP client with a conservative timeout for
// external provider API calls (Google userinfo, GitHub user/email). Using the
// default http.DefaultClient is dangerous because it has no timeout and can hold
// goroutines and DB transactions open indefinitely if the provider hangs.
var oauthHTTPClient = &http.Client{Timeout: 10 * time.Second}

// isSafeRedirectPath returns true when p is a relative path that is safe to
// use as a redirect target: it must start with "/" but not with "//" (which
// would create a protocol-relative URL pointing off-domain).
func isSafeRedirectPath(p string) bool {
	return p != "" && strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "//")
}

func NewAuthHandler(server *server.Server) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(server),
	}
}

func (h *AuthHandler) Login(c echo.Context) error {
	loginRequest := new(requests.LoginRequest)

	if err := c.Bind(loginRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := loginRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	// Use a single unified error for all login failures to prevent account enumeration.
	const loginErrMsg = "invalid email or password"

	user, err := h.Services(c).User.GetByEmail(c.Request().Context(), loginRequest.Email)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to get user by email", err)
		return responses.ErrorResponse(c, http.StatusUnauthorized, loginErrMsg)
	}

	authenticator, err := h.Services(c).Authenticator.GetAuthenticatorByType(c.Request().Context(), user.ID, models.AuthenticatorTypePassword)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to get authenticator by user id", err)
		return responses.ErrorResponse(c, http.StatusUnauthorized, loginErrMsg)
	}

	if !password_utils.ComparePasswords(authenticator.Token, []byte(loginRequest.Password)) {
		return responses.ErrorResponse(c, http.StatusUnauthorized, loginErrMsg)
	}

	session, err := h.Services(c).User.CreateUserSession(c.Request().Context(), user.ID, &user_service.CreateSessionParams{
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	})
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create session", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create session")
	}

	return responses.Response(c, http.StatusOK, responses.NewLoginResponse(session.Token, session.ExpiresAt.Format(time.RFC3339)))
}

func (h *AuthHandler) Signup(c echo.Context) error {
	signupRequest := new(requests.SignupRequest)

	if err := c.Bind(signupRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := signupRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	if h.Services(c).User.ExistsByEmail(c.Request().Context(), signupRequest.Email) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "User already exists")
	}

	user := models.User{
		Email: signupRequest.Email,
		Name:  signupRequest.Name,
	}

	err := h.Services(c).User.CreateUser(c.Request().Context(), &user)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create user", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create user")
	}

	hashedPassword := password_utils.GeneratePasswordHash([]byte(signupRequest.Password))

	authenticator := models.Authenticator{
		UserID: user.ID,
		Type:   models.AuthenticatorTypePassword,
		Token:  hashedPassword,
	}

	err = h.Services(c).Authenticator.CreateAuthenticator(c.Request().Context(), &authenticator)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create authenticator", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create authenticator")
	}

	h.Services(c).User.SendVerificationEmail(c.Request().Context(), &user)

	workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create default workspace", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create default workspace")
	}

	if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
		logger.Error(c.Request().Context(), "unable to set default workspace", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to set default workspace")
	}

	session, err := h.Services(c).User.CreateUserSession(c.Request().Context(), user.ID, &user_service.CreateSessionParams{
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	})
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create session", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create session")
	}

	return responses.Response(c, http.StatusOK, responses.NewLoginResponse(session.Token, session.ExpiresAt.Format(time.RFC3339)))
}

// ExchangeCode exchanges a short-lived OAuth exchange code for the real JWT session token.
// The code is a UUID issued by the OAuth callback; it is single-use and expires in 60 seconds.
func (h *AuthHandler) ExchangeCode(c echo.Context) error {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.Bind(&req); err != nil || req.Code == "" {
		return responses.ErrorResponse(c, http.StatusBadRequest, "code is required")
	}

	token, err := h.Services(c).Token.ValidateOAuthExchangeCode(c.Request().Context(), req.Code)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "invalid or expired code")
	}

	if err := h.Services(c).Token.MarkTokenAsUsed(c.Request().Context(), token.ID); err != nil {
		logger.Error(c.Request().Context(), "unable to mark exchange token as used", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "unable to complete sign in")
	}

	// token.Token is the real session JWT; token.ExpiresAt is the short-lived
	// exchange code expiry (60 s). Extract the actual session expiry from the JWT
	// so the frontend cookie is set with the correct lifetime.
	sessionExpiry, err := jwt_utils.GetJWTExpiry(token.Token)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to parse session expiry from token", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "unable to complete sign in")
	}
	return responses.Response(c, http.StatusOK, responses.NewLoginResponse(token.Token, sessionExpiry.Format(time.RFC3339)))
}

func (h *AuthHandler) GoogleAuthorize(c echo.Context) error {
	// authorize redirect for google oauth flow
	redirectURL := c.QueryParam("ref")
	state := config.NewOAuthState("google", config.OAuthFlowAuth, redirectURL, nil, nil, nil)
	state.SourcePage = c.QueryParam("page")
	url := config.OAuthProviders.Google.AuthCodeURL(state.ToString())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

// oauthErrorRedirect redirects the user to the source page (login or signup)
// with error query params. Falls back to /login when the source page is unknown.
func oauthErrorRedirect(c echo.Context, state *config.OAuthState, errorCode string) error {
	page := "login"
	if state != nil && state.SourcePage == "signup" {
		page = "signup"
	}
	redirectURL := config.Environment.ClientURL + "/" + page + "?error=true"
	if errorCode != "" {
		redirectURL += "&error_code=" + errorCode
	}
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) GoogleOAuthCallback(c echo.Context) error {
	// handle google oauth callback

	state := c.QueryParam("state")

	oauthState, err := config.ParseOAuthState(state)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to parse state from google", err)
		return oauthErrorRedirect(c, nil, "invalid_state")
	}

	code := c.QueryParam("code")
	token, err := config.OAuthProviders.Google.Exchange(c.Request().Context(), code)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to exchange code for token from google", err)
		return oauthErrorRedirect(c, oauthState, "")
	}
	encryptedRefreshToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.RefreshToken)
	if err != nil {
		return h.internalError(c, err)
	}
	encryptedAccessToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.AccessToken)
	if err != nil {
		return h.internalError(c, err)
	}

	googleReq, err := http.NewRequestWithContext(c.Request().Context(), http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return h.internalError(c, err)
	}
	googleReq.Header.Set("Authorization", "Bearer "+token.AccessToken)
	response, err := oauthHTTPClient.Do(googleReq)
	if err != nil {
		return oauthErrorRedirect(c, oauthState, "")
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to read response from google")
		return oauthErrorRedirect(c, oauthState, "")
	}

	var userInfo config.GoogleProfileResponse

	jsonUnmarshalError := json.Unmarshal(contents, &userInfo)
	if jsonUnmarshalError != nil {
		logger.Error(c.Request().Context(), "unable to parse response from google")
		return oauthErrorRedirect(c, oauthState, "")
	}

	/**
	* Steps to import user
		1. Check if a user exists with that email
			1.1. If found, check if the user has a google authenticator
					2.1 If found, update the access token for the user.
					2.2 If not found, create a new authenticator with the user
			1.2. If not found, create a user, and an authenticator with type Google
		2. Generate a session for the user
		3. Redirect the user to frontend with the generated session token
	*/

	exists := h.Services(c).User.ExistsByEmail(c.Request().Context(), userInfo.Email)

	var finalUser *models.User

	if !exists {
		var user models.User
		user.Email = userInfo.Email
		user.Name = userInfo.Name

		if userInfo.Picture != "" {
			user.Avatar = &userInfo.Picture
		}

		if user.VerifiedAt == nil {
			now := time.Now().UTC()
			user.VerifiedAt = &now
		}

		authenticator := models.Authenticator{
			UserID:       user.ID,
			Type:         models.AuthenticatorTypeGoogle,
			Token:        encryptedAccessToken,
			RefreshToken: &encryptedRefreshToken,
			Scopes:       strings.Join(config.OAuthProviders.Google.Scopes, ","),
		}

		err = h.Services(c).User.CreateUser(c.Request().Context(), &user)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create user", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_user")
		}

		authenticator.UserID = user.ID
		err = h.Services(c).Authenticator.CreateAuthenticator(c.Request().Context(), &authenticator)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create authenticator", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_authenticator")
		}

		workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create default workspace", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_workspace")
		}

		if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
			logger.Error(c.Request().Context(), "unable to set default workspace", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_set_default_workspace")
		}

		finalUser = &user
	} else {

		user, err := h.Services(c).User.GetByEmail(c.Request().Context(), userInfo.Email)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find user by email", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_find_user")
		}

		authenticators, err := h.Services(c).Authenticator.GetAuthenticatorsByUserID(c.Request().Context(), user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find authenticators by user id", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_find_authenticators")
		}

		var existingGoogleAuthenticator *models.Authenticator
		for _, authenticator := range authenticators {
			if authenticator.Type == models.AuthenticatorTypeGoogle {
				existingGoogleAuthenticator = &authenticator
			}
		}

		if existingGoogleAuthenticator != nil {
			err = h.Services(c).Authenticator.UpdateAuthenticatorTokens(c.Request().Context(), existingGoogleAuthenticator.ID, encryptedAccessToken, encryptedRefreshToken)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update authenticator", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_update_authenticator")
			}
		} else {
			authenticator := models.Authenticator{
				UserID:       user.ID,
				Type:         models.AuthenticatorTypeGoogle,
				Token:        encryptedAccessToken,
				RefreshToken: &encryptedRefreshToken,
				Scopes:       strings.Join(config.OAuthProviders.Google.Scopes, ","),
			}
			err = h.Services(c).Authenticator.CreateAuthenticator(c.Request().Context(), &authenticator)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to create authenticator", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_create_authenticator")
			}
		}

		if user.VerifiedAt == nil {
			err = h.Services(c).User.BulkVerifyUsers(c.Request().Context(), []uuid.UUID{user.ID})
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update user", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_update_user")
			}
		}

		finalUser = user
	}

	if finalUser == nil {
		logger.Error(c.Request().Context(), "Somehow, there's no user and we couldn't even create one")
		return oauthErrorRedirect(c, oauthState, "")
	}

	session, err := h.Services(c).User.CreateUserSession(c.Request().Context(), finalUser.ID, &user_service.CreateSessionParams{
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	})
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return oauthErrorRedirect(c, oauthState, "")
	}

	signupQueryParam := "false"
	if !exists {
		signupQueryParam = "true"
	}

	exchangeToken, err := h.Services(c).Token.CreateOAuthExchangeCode(c.Request().Context(), finalUser.ID, session.Token)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return oauthErrorRedirect(c, oauthState, "unable_to_create_exchange_code")
	}

	redirectURL := config.Environment.ClientURL + "/authenticated?code=" + exchangeToken.ID.String() + "&signup=" + signupQueryParam
	if oauthState.RedirectTo != "" && oauthState.RedirectTo != "undefined" && oauthState.RedirectTo != "null" {
		redirectURL += "&ref=" + url.QueryEscape(oauthState.RedirectTo)
	}
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func (h *AuthHandler) GitHubAuthorize(c echo.Context) error {
	redirectURL := c.QueryParam("ref")
	state := config.NewOAuthState("github", config.OAuthFlowAuth, redirectURL, nil, nil, nil)
	state.SourcePage = c.QueryParam("page")
	url := config.OAuthProviders.GitHubAuth.AuthCodeURL(state.ToString())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitHubOAuthCallback(c echo.Context) error {
	stateStr := c.QueryParam("state")

	oauthState, err := config.ParseOAuthState(stateStr)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to parse state from github", err)
		return oauthErrorRedirect(c, nil, "invalid_state")
	}

	code := c.QueryParam("code")
	token, err := config.OAuthProviders.GitHubAuth.Exchange(c.Request().Context(), code)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to exchange code for token from github", err)
		return oauthErrorRedirect(c, oauthState, "")
	}

	userProfile, err := fetchGitHubUserInfo(token.AccessToken)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to fetch github user info", err)
		return oauthErrorRedirect(c, oauthState, "")
	}

	// Integration-only flow: a logged-in user is connecting their GitHub account
	// to a workspace. Create/update the integration and redirect — do not log in.
	if oauthState.Flow == config.OAuthFlowIntegration {
		if oauthState.WorkspaceID == nil || oauthState.UserID == nil {
			logger.Error(c.Request().Context(), "integration oauth state missing workspace or user id")
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/integrations?error=invalid_state")
		}

		if upsertErr := h.upsertGitHubIntegration(c.Request().Context(), h.Services(c).Integration, *oauthState.WorkspaceID, *oauthState.UserID, token.AccessToken, &token.Expiry, userProfile); upsertErr != nil {
			logger.Error(c.Request().Context(), "unable to upsert github integration", upsertErr)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/integrations?error=creation_failed")
		}

		redirectPath := oauthState.RedirectTo
		if oauthState.CreateRepository != nil && *oauthState.CreateRepository {
			redirectPath = "/repositories/connect"
		}
		unescapedRedirectTo, err := url.QueryUnescape(redirectPath)
		if err != nil || !isSafeRedirectPath(unescapedRedirectTo) {
			unescapedRedirectTo = "/integrations"
		}
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+unescapedRedirectTo+"?success=true&type=github")
	}

	// Auth flow: login or sign up the user.
	encryptedAccessToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.AccessToken)
	if err != nil {
		return h.internalError(c, err)
	}

	// Resolve email — may be empty if the user has set it private.
	email := userProfile.Email
	if email == "" {
		email, err = fetchGitHubPrimaryEmail(token.AccessToken)
		if err != nil || email == "" {
			logger.Error(c.Request().Context(), "unable to fetch github user email", err)
			return oauthErrorRedirect(c, oauthState, "no_email")
		}
	}

	name := userProfile.Login
	if userProfile.Name != "" {
		name = userProfile.Name
	}

	exists := h.Services(c).User.ExistsByEmail(c.Request().Context(), email)

	var finalUser *models.User

	if !exists {
		var user models.User
		user.Email = email
		user.Name = name
		if userProfile.AvatarURL != "" {
			user.Avatar = &userProfile.AvatarURL
		}
		now := time.Now().UTC()
		user.VerifiedAt = &now

		authenticator := models.Authenticator{
			UserID: user.ID,
			Type:   models.AuthenticatorTypeGitHub,
			Token:  encryptedAccessToken,
			Scopes: strings.Join(config.OAuthProviders.GitHubAuth.Scopes, ","),
		}

		err = h.Services(c).User.CreateUser(c.Request().Context(), &user)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create user", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_user")
		}

		authenticator.UserID = user.ID
		err = h.Services(c).Authenticator.CreateAuthenticator(c.Request().Context(), &authenticator)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create authenticator", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_authenticator")
		}

		workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create default workspace", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_create_workspace")
		}

		if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
			logger.Error(c.Request().Context(), "unable to set default workspace", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_set_default_workspace")
		}

		// Auto-create the GitHub integration for the new workspace.
		if upsertErr := h.upsertGitHubIntegration(c.Request().Context(), h.Services(c).Integration, workspace.ID, user.ID, token.AccessToken, &token.Expiry, userProfile); upsertErr != nil {
			logger.Error(c.Request().Context(), "unable to auto-create github integration after signup", upsertErr)
		}

		finalUser = &user
	} else {
		user, err := h.Services(c).User.GetByEmail(c.Request().Context(), email)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find user by email", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_find_user")
		}

		authenticators, err := h.Services(c).Authenticator.GetAuthenticatorsByUserID(c.Request().Context(), user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find authenticators by user id", err)
			return oauthErrorRedirect(c, oauthState, "unable_to_find_authenticators")
		}

		var existingGitHubAuthenticator *models.Authenticator
		for _, authenticator := range authenticators {
			if authenticator.Type == models.AuthenticatorTypeGitHub {
				existingGitHubAuthenticator = &authenticator
				break
			}
		}

		if existingGitHubAuthenticator != nil {
			err = h.Services(c).Authenticator.UpdateAuthenticatorTokens(c.Request().Context(), existingGitHubAuthenticator.ID, encryptedAccessToken, "")
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update authenticator", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_update_authenticator")
			}
		} else {
			authenticator := models.Authenticator{
				UserID: user.ID,
				Type:   models.AuthenticatorTypeGitHub,
				Token:  encryptedAccessToken,
				Scopes: strings.Join(config.OAuthProviders.GitHubAuth.Scopes, ","),
			}
			err = h.Services(c).Authenticator.CreateAuthenticator(c.Request().Context(), &authenticator)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to create authenticator", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_create_authenticator")
			}
		}

		if user.VerifiedAt == nil {
			err = h.Services(c).User.BulkVerifyUsers(c.Request().Context(), []uuid.UUID{user.ID})
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update user", err)
				return oauthErrorRedirect(c, oauthState, "unable_to_update_user")
			}
		}

		finalUser = user
	}

	if finalUser == nil {
		logger.Error(c.Request().Context(), "no user found or created during github oauth")
		return oauthErrorRedirect(c, oauthState, "")
	}

	session, err := h.Services(c).User.CreateUserSession(c.Request().Context(), finalUser.ID, &user_service.CreateSessionParams{
		UserAgent: c.Request().UserAgent(),
		IP:        c.RealIP(),
	})
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return oauthErrorRedirect(c, oauthState, "")
	}

	signupQueryParam := "false"
	if !exists {
		signupQueryParam = "true"
	}

	exchangeToken, err := h.Services(c).Token.CreateOAuthExchangeCode(c.Request().Context(), finalUser.ID, session.Token)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return oauthErrorRedirect(c, oauthState, "unable_to_create_exchange_code")
	}

	redirectURL := config.Environment.ClientURL + "/authenticated?code=" + exchangeToken.ID.String() + "&signup=" + signupQueryParam
	if oauthState.RedirectTo != "" && oauthState.RedirectTo != "undefined" && oauthState.RedirectTo != "null" {
		redirectURL += "&ref=" + url.QueryEscape(oauthState.RedirectTo)
	}
	return c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// upsertGitHubIntegration creates or updates the GitHub integration for a workspace using the provided token.
func (h *AuthHandler) upsertGitHubIntegration(ctx context.Context, integrationSvc *integration_service.Service, workspaceID uuid.UUID, userID uuid.UUID, accessToken string, expiry *time.Time, profile *githubUserInfo) error {
	existing, err := integrationSvc.GetIntegrationByType(ctx, workspaceID, github_integration.IntegrationTypeGitHub)
	if err == nil && existing != nil {
		if tokenErr := integrationSvc.UpdateTokens(ctx, existing.ID, accessToken, nil, expiry, userID); tokenErr != nil {
			return tokenErr
		}
		existing.ServiceTitle = &profile.Login
		existing.ServiceIcon = &profile.AvatarURL
		existing.UpdatedBy = userID
		return integrationSvc.UpdateIntegration(ctx, existing)
	}

	login := profile.Login
	avatarURL := profile.AvatarURL
	idStr := strconv.Itoa(profile.ID)
	scopes := strings.Join(config.OAuthProviders.GitHubAuth.Scopes, ",")

	_, createErr := integrationSvc.CreateIntegration(ctx, &integration_service.CreateIntegrationRequest{
		WorkspaceID:     workspaceID,
		Type:            github_integration.IntegrationTypeGitHub,
		ServiceTitle:    &login,
		ServiceIcon:     &avatarURL,
		Token:           accessToken,
		TokenExpiresAt:  expiry,
		Scopes:          &scopes,
		ServiceTenantID: &idStr,
		ServiceUserID:   &login,
		CreatedBy:       userID,
	})
	return createErr
}

// fetchGitHubPrimaryEmail fetches the primary verified email for a GitHub user.
func fetchGitHubPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := oauthHTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var emails []config.GitHubEmailResponse
	if err := json.Unmarshal(body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}
	return "", nil
}

func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	verifyRequest := new(requests.VerifyEmailRequest)
	if err := c.Bind(verifyRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.Services(c).Token.GetTokenByValue(c.Request().Context(), verifyRequest.Token)
	if err != nil || token.Type != models.TokenTypeEmailVerification {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
	}

	err = h.Services(c).User.VerifyUser(c.Request().Context(), token.ID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to verify user")
	}

	return responses.Response(c, http.StatusOK, responses.NewGenericSuccessResponse(nil))
}

func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	user, err := context_utils.GetUserFromContext(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "user not found")
	}

	err = h.Services(c).User.SendVerificationEmail(c.Request().Context(), user)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to send verification email")
	}

	return responses.Response(c, http.StatusOK, responses.NewGenericSuccessResponse(nil))
}

func (h *AuthHandler) ResetPassword(c echo.Context) error {
	resetRequest := new(requests.ResetPasswordRequest)
	if err := c.Bind(resetRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	err := h.Services(c).User.SendResetPasswordEmail(c.Request().Context(), resetRequest.Email)
	if err != nil {
		logger.Error(c.Request().Context(), err)
	}

	return responses.Response(c, http.StatusOK, responses.NewGenericSuccessResponse(nil))
}

func (h *AuthHandler) ResetPasswordConfirm(c echo.Context) error {
	resetConfirmRequest := new(requests.ResetPasswordConfirmRequest)
	if err := c.Bind(resetConfirmRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.Services(c).Token.ValidateTokenByType(c.Request().Context(), resetConfirmRequest.Token, models.TokenTypePasswordReset)
	if err != nil || token.ID.String() != resetConfirmRequest.ID {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
	}

	err = h.Services(c).User.UpdatePassword(c.Request().Context(), token.UserID, resetConfirmRequest.Password)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to update password")
	}

	h.Services(c).Token.MarkTokenAsUsed(c.Request().Context(), token.ID)

	err = h.Services(c).Session.InvalidateAllSessionsForUser(c.Request().Context(), token.UserID)
	if err != nil {
		return h.internalError(c, err)
	}

	return responses.Response(c, http.StatusOK, responses.NewGenericSuccessResponse(nil))
}
