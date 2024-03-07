package v1

import (
	"encoding/json"

	"github.com/Woodfyn/file-api/internal/core"
	"github.com/gofiber/fiber"
)

func (h *Handler) initFileRouter(api fiber.Router) {
	h.authMiddleware()

	file := api.Group("/file")
	{
		file.Post("/upload", h.fileUpload)
		file.Get("/files", h.files)
	}
}

func (h *Handler) fileUpload(c *fiber.Ctx) {
	c.Set("Content-Type", "application/json")

	file, err := c.FormFile("fileUpload")
	if err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, "failed to get file")
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
		Reader: buffer,
	}

	if err := dto.Validate(); err != nil {
		newErrorResponse(c, fiber.StatusBadRequest, err.Error())
	}

	err = h.service.Files.Upload(c.Context(), &dto)
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "failed to upload file")
	}

	newSuccessResponse(c, fiber.StatusOK, "success")
}

func (h *Handler) files(c *fiber.Ctx) {
	c.Set("Content-Type", "application/json")

	files, err := h.service.Files.GetFiles(c.Context(), c.Fasthttp.Response.BodyWriter())
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "failed to get files")
	}

	resp := make([]map[string]interface{}, len(files))
	for i, fileInfo := range files {
		fileData := map[string]interface{}{
			"id":   fileInfo.ID,
			"name": fileInfo.Name,
			"size": fileInfo.Size,
			"file": fileInfo.Bytes,
		}
		resp[i] = fileData
	}

	respJSON, err := json.Marshal(resp)
	if err != nil {
		newErrorResponse(c, fiber.StatusInternalServerError, "error marshaling response")
		return
	}

	newSuccessResponse(c, fiber.StatusOK, string(respJSON))
}
