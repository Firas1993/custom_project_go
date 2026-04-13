package product

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler holds the product store and handles HTTP requests.
type Handler struct {
	store Store
}

// NewHandler returns a new Handler backed by the given store.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// Create handles POST /products.
func (h *Handler) Create(c *gin.Context) {
	var p Product
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	created, err := h.store.Create(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, created)
}

// List handles GET /products.
func (h *Handler) List(c *gin.Context) {
	products, err := h.store.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// Get handles GET /products/:id.
func (h *Handler) Get(c *gin.Context) {
	id := c.Param("id")
	p, err := h.store.Get(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Update handles PUT /products/:id.
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var p Product
	err := c.ShouldBindJSON(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.store.Update(id, p)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updated)
}

// Delete handles DELETE /products/:id.
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	err := h.store.Delete(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
