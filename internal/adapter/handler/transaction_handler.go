package handler

import (
	"net/http"

	_ "github.com/CardenalDex/crudprotec/internal/entitys"
	"github.com/CardenalDex/crudprotec/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	service usecase.TransactionUseCase
}

func NewTransactionHandler(s usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{service: s}
}

// Request DTO
type createTransactionRequest struct {
	MerchantID string  `json:"merchant_id" binding:"required"`
	Amount     float64 `json:"amount" binding:"required,gt=0"` // Input as float for user friendliness
}

// @Summary Create a new transaction
// @Description Calculates commission and fees, saves to DB, and returns the created transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body createTransactionRequest true "Transaction Request"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	merchantUUID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Merchant UUID"})
		return
	}

	// CONVERSION LAYER: Convert User Float ($200.00) -> System Int64 Cents (20000)
	amountCents := int64(req.Amount * 100)

	// Call Business Logic
	tx, err := h.service.ProcessTransaction(c.Request.Context(), merchantUUID, amountCents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}
