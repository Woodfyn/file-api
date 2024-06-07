package rest

import (
	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRouter(api fiber.Router) {
	auth := api.Group("/auth")
	{
		auth.Post("/sign-up", h.authSignUp)
		auth.Post("/sign-in", h.authSignIn)
		auth.Post("/refresh", h.authRefreshToken)
	}
}

func (h *Handler) authSignUp(c *fiber.Ctx) error {
	var input core.SignUpRequest
	if err := c.BodyParser(&input); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	if err := h.authService.SignUp(c.Context(), input); err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) authSignIn(c *fiber.Ctx) error {
	var input core.SignInRequest
	if err := c.BodyParser(&input); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	tokens, err := h.authService.SignIn(c.Context(), input)
	if err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}

	newDataResponse(c, fiber.StatusOK, tokens)
	return nil
}

func (h *Handler) authRefreshToken(c *fiber.Ctx) error {
	var input core.RefreshTokenReq
	if err := c.BodyParser(&input); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	if err := input.Validate(); err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	tokens, err := h.authService.Refresh(c.Context(), input.RefreshToken)
	if err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}

	newDataResponse(c, fiber.StatusOK, tokens)
	return nil
}
