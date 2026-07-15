package product

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Handler struct{ Repo *Repository }

func NewHandler(repo *Repository) *Handler { return &Handler{Repo: repo} }

func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad json"})
		return
	}
	req.SKU = strings.TrimSpace(req.SKU)
	req.Name = strings.TrimSpace(req.Name)
	if req.SKU == "" || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sku and name are required"})
		return
	}

	p, err := h.Repo.Create(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) List(c *gin.Context) {
	products, err := h.Repo.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
