package v1

import (
	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber"
)

func (h *Handler) initUserRouter(api fiber.Router) {
	user := api.Group("/auth")
	{
		user.Post("/sign-up", h.userSignUp)
		user.Post("/sing-in", h.userSignIn)
	}
}

func (h *Handler) userSignUp(c *fiber.Ctx) {
	var req core.SingUpRequest

	if err := c.BodyParser(&req); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// TODO: create user

	newSuccessResponse(c, fiber.StatusOK, "welcome")
}

func (h *Handler) userSignIn(c *fiber.Ctx) {
	var req core.SingInRequest

	if err := c.BodyParser(&req); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	if err := req.Validate(); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	// TODO: sing in user

	newSuccessResponse(c, fiber.StatusOK, "welcome")
}
