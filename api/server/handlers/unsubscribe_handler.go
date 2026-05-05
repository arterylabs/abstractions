package handlers

import (
	"abstractions/repositories"
	jwt_utils "abstractions/utils/jwt_utils"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// UnsubscribeHandler handles newsletter unsubscribe requests.
type UnsubscribeHandler struct {
	teamMemberRepo *repositories.TeamMemberRepository
}

func NewUnsubscribeHandler(db *gorm.DB) *UnsubscribeHandler {
	return &UnsubscribeHandler{
		teamMemberRepo: repositories.NewTeamMemberRepository(db),
	}
}

// Unsubscribe handles GET /unsubscribe?token={jwt}.
// The JWT is signed with the existing SIGNING_KEY and contains team_member_id.
// Sets subscribed=false for the identified team member and returns a confirmation page.
func (h *UnsubscribeHandler) Unsubscribe(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return c.String(http.StatusBadRequest, "missing token")
	}

	teamMemberID, err := jwt_utils.ParseUnsubscribeToken(token)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid or expired token")
	}

	ctx := c.Request().Context()
	if err := h.teamMemberRepo.UpdateSubscription(ctx, teamMemberID, false); err != nil {
		return c.String(http.StatusInternalServerError, "could not update subscription")
	}

	return c.HTML(http.StatusOK, `<!DOCTYPE html>
<html lang="en">
<head><meta charset="UTF-8"><title>Unsubscribed</title></head>
<body style="font-family:sans-serif;max-width:480px;margin:80px auto;text-align:center;">
<h1>You have been unsubscribed.</h1>
<p>You will no longer receive newsletters for this repository.</p>
</body>
</html>`)
}
