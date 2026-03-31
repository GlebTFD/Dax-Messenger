package http

import (
	"fmt"

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
		// add routing error depending on type
		h.log.Error("Error to delete msg", "error", err)
		c.Status(fiber.StatusInternalServerError)
		return fmt.Errorf("error to delete msg: %w", err)
	}

	return nil
}
