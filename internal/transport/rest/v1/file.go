package v1

import (
	"encoding/json"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber"
)

func (h *Handler) initFileRouter(api fiber.Router) {
	api.Use(h.authMiddleware())

	file := api.Group("/file")
	{
		file.Post("/upload", h.fileUpload)
		file.Get("/files", h.files)
	}
}

func (h *Handler) fileUpload(c *fiber.Ctx) {
	c.Accepts("multipart/form-data")

	file, err := c.FormFile("file")
	if err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, "invalid file")
		return
	}

	buffer, err := file.Open()
	if err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, "failed to open file")
		return
	}

	defer buffer.Close()

	dto := core.CreateFileDTO{
		Name:   file.Filename,
		Size:   file.Size,
		Type:   file.Header.Get("Content-Type"),
		Reader: buffer,
	}

	if err := h.service.Files.Upload(c.Context(), &dto); err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "failed to upload file")
		return
	}

	newSuccessResponse(c, fiber.StatusOK, "success")
}

func (h *Handler) files(c *fiber.Ctx) {
	files, err := h.service.Files.GetFiles(c.Context())
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "failed to get files")
	}

	response := make([]map[string]interface{}, len(files))
	for i, fileInfo := range files {
		fileData := map[string]interface{}{
			"id":   fileInfo.ID,
			"name": fileInfo.Name,
			"size": fileInfo.Size,
			"file": fileInfo.Bytes,
		}
		response[i] = fileData
	}

	respJSON, err := json.Marshal(response)
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "error marshaling response")
		return
	}

	newSuccessResponse(c, fiber.StatusOK, string(respJSON))
}
