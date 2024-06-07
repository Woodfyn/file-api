package rest

import (
	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initFileRouter(api fiber.Router) {
	api.Use(h.authMiddleware())

	file := api.Group("/file")
	{
		file.Use(h.authMiddleware())

		file.Post("/upload", h.fileUpload)
		file.Get("/", h.fileGetAll)
	}
}

func (h *Handler) fileUpload(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return newErrorResponse(fiber.StatusUnauthorized, core.ErrInvalidAccessToken)
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return newErrorResponse(fiber.StatusBadRequest, err)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}
	defer file.Close()

	if err := h.fileService.Upload(c.Context(), fileHeader, file, userId); err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}

	return c.SendStatus(fiber.StatusCreated)
}

func (h *Handler) fileGetAll(c *fiber.Ctx) error {
	userId, ok := c.Locals("userId").(string)
	if !ok {
		return newErrorResponse(fiber.StatusUnauthorized, core.ErrInvalidAccessToken)
	}

	response, err := h.fileService.GetFiles(c.Context(), userId)
	if err != nil {
		return newErrorResponse(fiber.StatusInternalServerError, err)
	}

	newDataResponse(c, fiber.StatusOK, response)
	return nil
}
