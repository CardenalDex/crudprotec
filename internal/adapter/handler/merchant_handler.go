package handler

import (
	"net/http"

	_ "github.com/CardenalDex/crudprotec/internal/entitys" // needed for swager
	"github.com/CardenalDex/crudprotec/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MerchantHandler struct {
	service usecase.MerchantUseCase
}

func NewMerchantHandler(s usecase.MerchantUseCase) *MerchantHandler {
	return &MerchantHandler{service: s}
}

type createMerchantRequest struct {
	BusinessID string `json:"business_id" binding:"required"`
}

// --- Handlers ---

// @Summary Register a new Merchant
// @Description Creates a new merchant linked to a specific business
// @Tags merchants
// @Accept json
// @Produce json
// @Param merchant body createMerchantRequest true "Merchant Registration"
// @Success 201 {object} entity.Merchant
// @Failure 400 {object} map[string]string "Invalid input or Business UUID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /merchants/new [post]
func (h *MerchantHandler) RegisterMerchant(c *gin.Context) {
	var req createMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bizUUID, err := uuid.Parse(req.BusinessID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Business UUID"})
		return
	}

	merchant, err := h.service.RegisterMerchant(c.Request.Context(), bizUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, merchant)
}

// @Summary Get Merchant Details
// @Description Retrieve details of a specific merchant by ID
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant UUID"
// @Success 200 {object} entity.Merchant
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 404 {object} map[string]string "Merchant not found"
// @Router /merchants/{id} [get]
func (h *MerchantHandler) GetMerchant(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	merchant, err := h.service.GetMerchant(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchant)
}

// @Summary List Merchants by Business
// @Description Retrieve all merchants belonging to a specific Business
// @Tags merchants
// @Produce json
// @Param businessID path string true "Business UUID"
// @Success 200 {array} entity.Merchant
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /merchants/bybusiness/{businessID} [get]
func (h *MerchantHandler) GetBusinessMerchants(c *gin.Context) {
	idParam := c.Param("businessID")
	bizID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Business UUID format"})
		return
	}

	merchants, err := h.service.GetBusinessMerchants(c.Request.Context(), bizID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, merchants)
}

// @Summary Remove a Merchant
// @Description Logic delete of a merchant
// @Tags merchants
// @Produce json
// @Param id path string true "Merchant UUID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /merchants/delete/{id} [delete]
func (h *MerchantHandler) RemoveMerchant(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	if err := h.service.RemoveMerchant(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Merchant removed successfully"})
}
