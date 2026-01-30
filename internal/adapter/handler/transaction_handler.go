package handler

import (
	"net/http"

	_ "github.com/CardenalDex/crudprotec/internal/entitys" // neded for swagger
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
	Amount     float64 `json:"amount" binding:"required,gt=0"` // Input as float for user friendliness (converted after)
}

// @Summary Create a new transaction
// @Description Calculates commission and fees, saves to DB, and returns the created transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Param actor header string false "The name of the user performing the action"
// @Param transaction body createTransactionRequest true "Transaction Request"
// @Success 201 {object} entity.Transaction
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/new [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req createTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actor := c.GetHeader("actor")
	merchantUUID, err := uuid.Parse(req.MerchantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Merchant UUID"})
		return
	}

	// CONVERSION LAYER: Convert User Float ($200.00) -> System Int64 Cents (20000)
	amountCents := int64(req.Amount * 100)

	tx, err := h.service.ProcessTransaction(c.Request.Context(), actor, merchantUUID, amountCents)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tx)
}

// @Summary Get a transaction by ID
// @Description Retrieve a specific transaction details by its UUID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction UUID"
// @Success 200 {object} entity.Transaction
// @Failure 400 {object} map[string]string "Invalid UUID format"
// @Failure 404 {object} map[string]string "Transaction not found"
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	idParam := c.Param("id")
	txID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID format"})
		return
	}

	tx, err := h.service.GetTransaction(c.Request.Context(), txID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, tx)
}

// @Summary List transactions by Merchant
// @Description Retrieve all transactions belonging to a specific merchant
// @Tags transactions
// @Produce json
// @Param merchantID path string true "Merchant UUID"
// @Success 200 {array} entity.Transaction
// @Failure 400 {object} map[string]string "Invalid UUID format"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/bymerchant/{merchantID} [get]
func (h *TransactionHandler) GetMerchantTransactions(c *gin.Context) {
	mParam := c.Param("merchantID")
	mID, err := uuid.Parse(mParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Merchant UUID format"})
		return
	}

	transactions, err := h.service.GetMerchantTransactions(c.Request.Context(), mID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary List all transactions
// @Description Retrieve all transactions in the system
// @Tags transactions
// @Produce json
// @Success 200 {array} entity.Transaction
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/transactions [get]
func (h *TransactionHandler) GetAllTransactions(c *gin.Context) {
	transactions, err := h.service.GetAllTransactions(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary Get all revenue
// @Description Retrieve the total revenue of the system
// @Tags transactions
// @Produce json
// @Success 200 {number} number "Example: 1.25"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/revenue [get]
func (h *TransactionHandler) GetAllRevenue(c *gin.Context) {
	transactions, err := h.service.GetAllRevenue(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// @Summary Get revenue by Merchant
// @Description Retrieve all revenue of a merchant
// @Tags transactions
// @Produce json
// @Param merchantID path string true "Merchant UUID"
// @Success 200 {number} number "Example: 1.25"
// @Failure 400 {object} map[string]string "Invalid UUID format"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /transactions/revenuebymerchant/{merchantID} [get]
func (h *TransactionHandler) GetAllRevenueByMerchant(c *gin.Context) {
	mParam := c.Param("merchantID")
	mID, err := uuid.Parse(mParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Merchant UUID format"})
		return
	}

	transactions, err := h.service.GetAllRevenueByMerchant(c.Request.Context(), mID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}
