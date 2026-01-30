package handler

import (
	"net/http"

	_ "github.com/CardenalDex/crudprotec/internal/entitys" // neeeded for swagger
	"github.com/CardenalDex/crudprotec/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	service usecase.AdminUseCase
}

func NewAdminHandler(s usecase.AdminUseCase) *AdminHandler {
	return &AdminHandler{service: s}
}

type createBusinessRequest struct {
	Commission float64 `json:"commission_percentage" binding:"required,gt=0"` // e.g., 5.5 for 5.5%
}

type updateCommissionRequest struct {
	Commission float64 `json:"new_commission_percentage" binding:"required,gt=0"`
}

// --- Handlers ---

// @Summary Register a new Business
// @Description Creates a new business entity with a specific commission rate
// @Tags admin
// @Accept json
// @Produce json
// @Param actor header string false "The name of the user performing the action"
// @Param business body createBusinessRequest true "Business Registration"
// @Success 201 {object} entity.Business
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/businesses/new [post]
func (h *AdminHandler) RegisterBusiness(c *gin.Context) {
	var req createBusinessRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actor := c.GetHeader("actor")
	// Convert Percentage (5.5) -> Basis Points (550)
	commissionBP := int64(req.Commission * 100)

	biz, err := h.service.RegisterBusiness(c.Request.Context(), actor, commissionBP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, biz)
}

// @Summary Get Business Details
// @Description Retrieve details of a specific business by ID
// @Tags admin
// @Produce json
// @Param id path string true "Business UUID"
// @Success 200 {object} entity.Business
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 404 {object} map[string]string "Business not found"
// @Router /admin/businesses/{id} [get]
func (h *AdminHandler) GetBusiness(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	biz, err := h.service.GetBusiness(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, biz)
}

// @Summary Update Business Commission
// @Description Update the commission rate for an existing business
// @Tags admin
// @Accept json
// @Produce json
// @Param actor header string false "The name of the user performing the action"
// @Param id path string true "Business UUID"
// @Param commission body updateCommissionRequest true "New Commission Rate"
// @Success 200 {object} entity.Business
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/businesses/{id}/commission [patch]
func (h *AdminHandler) UpdateBusinessCommission(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	actor := c.GetHeader("actor")
	var req updateCommissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert Percentage -> Basis Points
	commissionBP := int64(req.Commission * 100)

	biz, err := h.service.UpdateBusinessCommission(c.Request.Context(), actor, id, commissionBP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, biz)
}

// @Summary Delete a Business
// @Description Soft delete a business (Logical Delete)
// @Tags admin
// @Produce json
// @Param actor header string false "The name of the user performing the action"
// @Param id path string true "Business UUID"
// @Success 200 {object} map[string]string "Success message"
// @Failure 400 {object} map[string]string "Invalid UUID"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /admin/businesses/delete/{id} [delete]
func (h *AdminHandler) RemoveBusiness(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}
	actor := c.GetHeader("actor")
	if err := h.service.RemoveBusiness(c.Request.Context(), actor, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Business deleted successfully"})
}

// @Summary Get Audit Logs
// @Description Retrieve audit logs for a specific resource (e.g., a Business ID or Transaction ID)
// @Tags audit
// @Produce json
// @Param resource_id path string true "Resource ID"
// @Success 200 {array} entity.Log
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /audit/{resource_id} [get]
func (h *AdminHandler) GetAuditTrail(c *gin.Context) {
	resourceID := c.Param("resource_id")

	logs, err := h.service.GetAuditTrail(c.Request.Context(), resourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

// @Summary Get All Logs
// @Description Retrieve all system audit logs
// @Tags audit
// @Produce json
// @Success 200 {array} entity.Log
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /audit [get]
func (h *AdminHandler) GetAllLogs(c *gin.Context) {
	logs, err := h.service.GetAllLogs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
