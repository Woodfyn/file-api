package v1

import (
	"errors"
	"strings"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber"
)

func (h *Handler) initUserRouter(api fiber.Router) {
	user := api.Group("/auth")
	{
		user.Post("/sign-up", h.userSignUp)
		user.Post("/sign-in", h.userSignIn)
		user.Get("/refresh", h.userRefreshToken)
	}
}

func (h *Handler) userSignUp(c *fiber.Ctx) {
	c.Accepts("application/json")

	var req core.SignUpRequest

	if err := c.BodyParser(&req); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Users.SignUp(c.Context(), req); err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		return
	}

	newSuccessResponse(c, fiber.StatusOK, "success")
}

func (h *Handler) userSignIn(c *fiber.Ctx) {
	c.Accepts("application/json")

	var req core.SignInRequest

	if err := c.BodyParser(&req); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
		return
	}

	accessToken, refreshToken, err := h.service.Users.SignIn(c.Context(), req)
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "Bearer " + refreshToken,
		HTTPOnly: true,
		Path:     "/",
	})

	newSuccessResponse(c, fiber.StatusOK, accessToken)
}

func (h *Handler) userRefreshToken(c *fiber.Ctx) {
	c.Accepts("application/json")

	var req core.Token

	if err := c.BodyParser(&req); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	cookie := c.Cookies("Authorization")
	refreshToken, err := getTokenFromCookie(cookie)
	if err != nil {
		newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		return
	}

	accessToken, err := h.service.Users.Refresh(c.Context(), refreshToken)
	if errors.Is(err, core.ErrTokenExpired) {
		newErrorResponse(c, fiber.StatusUnauthorized, err.Error())
		return
	} else if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, err.Error())
		return
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    "Bearer " + refreshToken,
		HTTPOnly: true,
		Path:     "/",
	})

	newSuccessResponse(c, fiber.StatusOK, accessToken)
}

func getTokenFromCookie(cookieValue string) (string, error) {
	headerParts := strings.Split(cookieValue, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
