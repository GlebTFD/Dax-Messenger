package http

import (
	"errors"

	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/dto"
	"github.com/GlebTFD/Dax-Messenger/Messenge-service/internal/usecase"
	"github.com/gofiber/fiber/v3"
	"github.com/hashicorp/go-hclog"
)

type HTTPHandeler struct {
	log     hclog.Logger
	service *usecase.MessageService
}

func NewHTTPHandler(log hclog.Logger, service *usecase.MessageService) *HTTPHandeler {
	return &HTTPHandeler{
		log:     log,
		service: service,
	}
}

func (h *HTTPHandeler) DeleteMessage(c fiber.Ctx) error {
	id := c.Params("id")

	err := h.service.DeleteMessage(id)
	if err != nil {
		if errors.Is(err, usecase.ErrMessageNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.SendString(usecase.ErrMessageNotFound.Error())
		}
		h.log.Error("Error to delete msg", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.SendString("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.DeleteMessageResponse{
		ID:      id,
		Deleted: true,
	})
}

func (h *HTTPHandeler) UpdateMessage(c fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateMessageRequest
	if err := c.Bind().JSON(&req); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.SendString("invalid request body")
	}

	err := h.service.UpdateMessage(id, req.Text)
	if err != nil {
		if errors.Is(err, usecase.ErrEmptyText) {
			c.Status(fiber.StatusBadRequest)
			return c.SendString(usecase.ErrEmptyText.Error())
		}
		if errors.Is(err, usecase.ErrMessageNotFound) {
			c.Status(fiber.StatusNotFound)
			return c.SendString(usecase.ErrMessageNotFound.Error())
		}
		h.log.Error("Error to update msg", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return c.SendString("internal server error")
	}

	return c.Status(fiber.StatusOK).JSON(dto.UpdateMessageResponse{
		ID:      id,
		Updated: true,
		Text:    req.Text,
	})
}
