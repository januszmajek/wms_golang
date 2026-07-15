package stock

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct{ Service *Service }

func NewHandler(service *Service) *Handler { return &Handler{Service: service} }

func (h *Handler) Inbound(c *gin.Context) {
	var req InboundRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	resp, err := h.Service.Receive(req)
	if errors.Is(err, ErrBadQuantity) || errors.Is(err, ErrProductMissing) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, resp)
}

func (h *Handler) Report(c *gin.Context) {
	items, err := h.Service.Report()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}
