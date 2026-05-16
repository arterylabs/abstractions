package handlers

import (
	"abstractions/config"
	"abstractions/db/models"
	"abstractions/logger"
	"abstractions/server"
	"abstractions/server/requests"
	"abstractions/server/responses"
	"abstractions/utils/context_utils"
	"abstractions/utils/encryption_utils"
	"abstractions/utils/jwt_utils"
	"abstractions/utils/password_utils"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthHandler struct {
	*BaseHandler
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

	user, err := h.Services(c).User.GetByEmail(loginRequest.Email)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to get user by email", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "User not found")
	}

	authenticator, err := h.Services(c).Authenticator.GetAuthenticatorByType(user.ID, models.AuthenticatorTypePassword)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to get authenticator by user id", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "Authenticator not found")
	}

	if !password_utils.ComparePasswords(authenticator.Token, []byte(loginRequest.Password)) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid password")
	}

	expires_at := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))

	sessionID := uuid.New()
	jwtToken, err := jwt_utils.GenerateJWT(user.ID.String(), sessionID.String(), expires_at)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to generate jwt token", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to generate jwt token")
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		UserAgent: c.Request().UserAgent(),
		IP:        c.Request().RemoteAddr,
		Token:     jwtToken,
		ExpiresAt: expires_at,
	}

	err = h.Services(c).Session.CreateSession(&session)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create session", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create session")
	}

	return responses.Response(c, http.StatusOK, responses.NewLoginResponse(session.Token, expires_at.Format(time.RFC3339)))
}

func (h *AuthHandler) Signup(c echo.Context) error {
	signupRequest := new(requests.SignupRequest)

	if err := c.Bind(signupRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	if err := signupRequest.Validate(); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid")
	}

	if h.Services(c).User.ExistsByEmail(signupRequest.Email) {
		return responses.ErrorResponse(c, http.StatusBadRequest, "User already exists")
	}

	user := models.User{
		Email: signupRequest.Email,
		Name:  signupRequest.Name,
	}

	err := h.Services(c).User.CreateUser(&user)
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

	err = h.Services(c).Authenticator.CreateAuthenticator(&authenticator)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create authenticator", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create authenticator")
	}

	h.Services(c).User.SendVerificationEmail(&user)

	workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create default workspace", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create default workspace")
	}

	if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
		logger.Error(c.Request().Context(), "unable to set default workspace", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to set default workspace")
	}

	expires_at := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))

	sessionID := uuid.New()
	jwtToken, err := jwt_utils.GenerateJWT(user.ID.String(), sessionID.String(), expires_at)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to generate jwt token", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to generate jwt token")
	}

	session := models.Session{
		ID:        sessionID,
		UserID:    user.ID,
		UserAgent: c.Request().UserAgent(),
		IP:        c.Request().RemoteAddr,
		Token:     jwtToken,
		ExpiresAt: expires_at,
	}

	err = h.Services(c).Session.CreateSession(&session)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to create session", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to create session")
	}

	return responses.Response(c, http.StatusOK, responses.NewSignupResponse(session.Token, expires_at.Format(time.RFC3339)))
}

func (h *AuthHandler) GoogleAuthorize(c echo.Context) error {
	// authorize redirect for google oauth flow
	redirectURL := c.QueryParam("ref")
	// stateValue := "abstractions"
	// if redirectURL != "" {
	// 	stateValue = stateValue + "<__>" + redirectURL
	// }

	state := config.NewOAuthState("google", redirectURL, nil, nil, nil)
	url := config.OAuthProviders.Google.AuthCodeURL(state.ToString())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c echo.Context) error {
	// handle google oauth callback

	state := c.QueryParam("state")

	oauthState, err := config.ParseOAuthState(state)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to parse state from google", err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=invalid_state")
	}

	code := c.QueryParam("code")
	token, err := config.OAuthProviders.Google.Exchange(c.Request().Context(), code)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to exchange code for token from google", err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}
	encryptedRefreshToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.RefreshToken)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	encryptedAccessToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.AccessToken)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}
	defer response.Body.Close()
	contents, err := io.ReadAll(response.Body)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to read response from google")
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	var userInfo config.GoogleProfileResponse

	jsonUnmarshalError := json.Unmarshal(contents, &userInfo)
	if jsonUnmarshalError != nil {
		logger.Error(c.Request().Context(), "unable to parse response from google")
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
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

	exists := h.Services(c).User.ExistsByEmail(userInfo.Email)

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

		err = h.Services(c).User.CreateUser(&user)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create user", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_user")
		}

		authenticator.UserID = user.ID
		err = h.Services(c).Authenticator.CreateAuthenticator(&authenticator)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create authenticator", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_authenticator")
		}

		workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create default workspace", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_workspace")
		}

		if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
			logger.Error(c.Request().Context(), "unable to set default workspace", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_set_default_workspace")
		}

		finalUser = &user
	} else {

		user, err := h.Services(c).User.GetByEmail(userInfo.Email)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find user by email", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_find_user")
		}

		authenticators, err := h.Services(c).Authenticator.GetAuthenticatorsByUserID(user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find authenticators by user id", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_find_authenticators")
		}

		var existingGoogleAuthenticator *models.Authenticator
		for _, authenticator := range authenticators {
			if authenticator.Type == "google" {
				existingGoogleAuthenticator = &authenticator
			}
		}

		if existingGoogleAuthenticator != nil {
			err = h.Services(c).Authenticator.UpdateAuthenticatorTokens(existingGoogleAuthenticator.ID, encryptedAccessToken, encryptedRefreshToken)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update authenticator", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_update_authenticator")
			}
		} else {
			authenticator := models.Authenticator{
				UserID:       user.ID,
				Type:         models.AuthenticatorTypeGoogle,
				Token:        encryptedAccessToken,
				RefreshToken: &encryptedRefreshToken,
				Scopes:       strings.Join(config.OAuthProviders.Google.Scopes, ","),
			}
			err = h.Services(c).Authenticator.CreateAuthenticator(&authenticator)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to create authenticator", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_authenticator")
			}
		}

		if user.VerifiedAt == nil {
			err = h.Services(c).User.BulkVerifyUsers([]uuid.UUID{user.ID})
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update user", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_update_user")
			}
		}

		finalUser = user
	}

	if finalUser == nil {
		logger.Error(c.Request().Context(), "Somehow, there's no user and we couldn't even create one")
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	expires_at := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))

	sessionID := uuid.New()

	session := models.Session{
		ID:        sessionID,
		UserID:    finalUser.ID,
		UserAgent: c.Request().UserAgent(),
		IP:        c.Request().RemoteAddr,
		ExpiresAt: expires_at,
	}

	sessionToken, err := jwt_utils.GenerateJWT(finalUser.ID.String(), session.ID.String(), expires_at)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	session.Token = sessionToken

	err = h.Services(c).Session.CreateSession(&session)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	signupQueryParam := "false"
	if !exists {
		signupQueryParam = "true"
	}

	if oauthState.RedirectTo != "" && oauthState.RedirectTo != "undefined" && oauthState.RedirectTo != "null" {
		// If the redirect URL is provided, redirect to that URL with the session token and user details
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/authenticated?token="+session.Token+"&expiresAt="+session.ExpiresAt.String()+"&email="+finalUser.Email+"&uid="+finalUser.ID.String()+"&ref="+oauthState.RedirectTo+"&signup="+signupQueryParam)
	}

	return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/authenticated?token="+session.Token+"&expiresAt="+session.ExpiresAt.String()+"&email="+finalUser.Email+"&uid="+finalUser.ID.String()+"&signup="+signupQueryParam)
}

func (h *AuthHandler) GitHubAuthorize(c echo.Context) error {
	redirectURL := c.QueryParam("ref")
	state := config.NewOAuthState("github", redirectURL, nil, nil, nil)
	url := config.OAuthProviders.GitHubAuth.AuthCodeURL(state.ToString())
	return c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GitHubCallback(c echo.Context) error {
	state := c.QueryParam("state")

	oauthState, err := config.ParseOAuthState(state)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to parse state from github", err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=invalid_state")
	}

	code := c.QueryParam("code")
	token, err := config.OAuthProviders.GitHubAuth.Exchange(c.Request().Context(), code)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to exchange code for token from github", err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	encryptedAccessToken, err := encryption_utils.EncryptString(config.Environment.EncryptionKey32Bit, token.AccessToken)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	// Fetch user profile
	userProfile, err := fetchGitHubUserInfo(token.AccessToken)
	if err != nil {
		logger.Error(c.Request().Context(), "unable to fetch github user info", err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	// Resolve email — may be empty if the user has set it private
	email := userProfile.Email
	if email == "" {
		email, err = fetchGitHubPrimaryEmail(token.AccessToken)
		if err != nil || email == "" {
			logger.Error(c.Request().Context(), "unable to fetch github user email", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=no_email")
		}
	}

	name := userProfile.Login
	if userProfile.Name != "" {
		name = userProfile.Name
	}

	exists := h.Services(c).User.ExistsByEmail(email)

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

		err = h.Services(c).User.CreateUser(&user)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create user", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_user")
		}

		authenticator.UserID = user.ID
		err = h.Services(c).Authenticator.CreateAuthenticator(&authenticator)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create authenticator", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_authenticator")
		}

		workspace, err := h.Services(c).Workspace.CreateWorkspace(c.Request().Context(), "My Workspace", user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to create default workspace", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_workspace")
		}

		if err := h.Services(c).User.UpdateDefaultWorkspace(c.Request().Context(), user.ID, workspace.ID); err != nil {
			logger.Error(c.Request().Context(), "unable to set default workspace", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_set_default_workspace")
		}

		finalUser = &user
	} else {
		user, err := h.Services(c).User.GetByEmail(email)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find user by email", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_find_user")
		}

		authenticators, err := h.Services(c).Authenticator.GetAuthenticatorsByUserID(user.ID)
		if err != nil {
			logger.Error(c.Request().Context(), "unable to find authenticators by user id", err)
			return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_find_authenticators")
		}

		var existingGitHubAuthenticator *models.Authenticator
		for _, authenticator := range authenticators {
			if authenticator.Type == models.AuthenticatorTypeGitHub {
				existingGitHubAuthenticator = &authenticator
				break
			}
		}

		if existingGitHubAuthenticator != nil {
			err = h.Services(c).Authenticator.UpdateAuthenticatorTokens(existingGitHubAuthenticator.ID, encryptedAccessToken, "")
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update authenticator", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_update_authenticator")
			}
		} else {
			authenticator := models.Authenticator{
				UserID: user.ID,
				Type:   models.AuthenticatorTypeGitHub,
				Token:  encryptedAccessToken,
				Scopes: strings.Join(config.OAuthProviders.GitHubAuth.Scopes, ","),
			}
			err = h.Services(c).Authenticator.CreateAuthenticator(&authenticator)
			if err != nil {
				logger.Error(c.Request().Context(), "unable to create authenticator", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_create_authenticator")
			}
		}

		if user.VerifiedAt == nil {
			err = h.Services(c).User.BulkVerifyUsers([]uuid.UUID{user.ID})
			if err != nil {
				logger.Error(c.Request().Context(), "unable to update user", err)
				return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true&error_code=unable_to_update_user")
			}
		}

		finalUser = user
	}

	if finalUser == nil {
		logger.Error(c.Request().Context(), "no user found or created during github oauth")
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	expires_at := time.Now().UTC().Add(time.Hour * 24 * 7 * time.Duration(6))
	sessionID := uuid.New()

	session := models.Session{
		ID:        sessionID,
		UserID:    finalUser.ID,
		UserAgent: c.Request().UserAgent(),
		IP:        c.Request().RemoteAddr,
		ExpiresAt: expires_at,
	}

	sessionToken, err := jwt_utils.GenerateJWT(finalUser.ID.String(), session.ID.String(), expires_at)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	session.Token = sessionToken

	err = h.Services(c).Session.CreateSession(&session)
	if err != nil {
		logger.Error(c.Request().Context(), err)
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/login?error=true")
	}

	signupQueryParam := "false"
	if !exists {
		signupQueryParam = "true"
	}

	if oauthState.RedirectTo != "" && oauthState.RedirectTo != "undefined" && oauthState.RedirectTo != "null" {
		return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/authenticated?token="+session.Token+"&expiresAt="+session.ExpiresAt.String()+"&email="+finalUser.Email+"&uid="+finalUser.ID.String()+"&ref="+oauthState.RedirectTo+"&signup="+signupQueryParam)
	}

	return c.Redirect(http.StatusTemporaryRedirect, config.Environment.ClientURL+"/authenticated?token="+session.Token+"&expiresAt="+session.ExpiresAt.String()+"&email="+finalUser.Email+"&uid="+finalUser.ID.String()+"&signup="+signupQueryParam)
}

// fetchGitHubPrimaryEmail fetches the primary verified email for a GitHub user.
func fetchGitHubPrimaryEmail(accessToken string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := http.DefaultClient.Do(req)
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

	token, err := h.Services(c).Token.GetTokenByValue(verifyRequest.Token)
	if err != nil || token.Type != models.TokenTypeEmailVerification {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
	}

	err = h.Services(c).User.VerifyUser(token.ID)
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

	err = h.Services(c).User.SendVerificationEmail(user)
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

	err := h.Services(c).User.SendResetPasswordEmail(resetRequest.Email)
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

	token, err := h.Services(c).Token.ValidateTokenByType(resetConfirmRequest.Token, models.TokenTypePasswordReset)
	if err != nil || token.ID.String() != resetConfirmRequest.ID {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
	}

	err = h.Services(c).User.UpdatePassword(token.UserID, resetConfirmRequest.Password)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to update password")
	}

	h.Services(c).Token.MarkTokenAsUsed(token.ID)

	err = h.Services(c).Session.InvalidateAllSessionsForUser(token.UserID)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return responses.Response(c, http.StatusOK, responses.NewGenericSuccessResponse(nil))
}
