package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ditto/internal/model"
	"ditto/internal/service"
)

// ThingHandler handles HTTP requests for things
type ThingHandler struct {
	service *service.ThingService
}

// NewThingHandler creates a new thing handler
func NewThingHandler(service *service.ThingService) *ThingHandler {
	return &ThingHandler{
		service: service,
	}
}

// Create handles the creation of a new thing
func (h *ThingHandler) Create(c *gin.Context) {
	var input model.ThingCreate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thing, err := h.service.Create(c.Request.Context(), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, thing)
}

// GetByID handles retrieving a thing by its ID
func (h *ThingHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	thing, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "thing not found"})
		return
	}

	c.JSON(http.StatusOK, thing)
}

// Update handles updating an existing thing
func (h *ThingHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var input model.ThingUpdate
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thing, err := h.service.Update(c.Request.Context(), id, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thing)
}

// Delete handles the deletion of a thing
func (h *ThingHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// List handles retrieving a list of things with pagination
func (h *ThingHandler) List(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	things, err := h.service.List(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, things)
}
