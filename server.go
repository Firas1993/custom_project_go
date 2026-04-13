package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mime-rona/irp-app-from-template/internal/product"
)

// CreateRouter sets up and returns the application HTTP router.
func CreateRouter() *gin.Engine {
	router := gin.Default()
	store := product.NewMemoryStore()
	h := product.NewHandler(store)

	products := router.Group("/products")
	products.POST("", h.Create)
	products.GET("", h.List)
	products.GET("/:id", h.Get)
	products.PUT("/:id", h.Update)
	products.DELETE("/:id", h.Delete)

	return router
}
