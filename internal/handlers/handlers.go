package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"FinalTaskAppGoBasic/internal/currencies"
	"FinalTaskAppGoBasic/internal/logs"
	"FinalTaskAppGoBasic/internal/models"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type Handlers struct {
	gorm *gorm.DB
}

// HandleTransactions роутер для обработки запросов к транзакциям
func (h *Handlers) HandleTransactions(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		currency := r.URL.Query().Get("currency")
		if currency != "" {
			h.getTransactionsWithCurrency(w, r, currency)

			return
		}

		h.getTransactions(w, r)
	case "POST":
		h.addTransaction(w, r)
	case "PUT":
		h.updateTransaction(w, r)
	case "DELETE":
		h.deleteTransaction(w, r)
	default:
		logs.HandleMessage(w, r, http.StatusNotFound, "method not supported")
	}
}

func (h *Handlers) addTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := h.parseTokenFromRequest(r)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "parse token failed")

		return
	}

	if userID == "" {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "user id is required")

		return
	}

	var tx models.Transactions
	err = json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	h.calculateCommission(&tx)

	tx.UserID = userID
	tx.TransactionDate = time.Now()

	result := h.gorm.Create(&tx)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "add transaction failed")

		return

	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logs.Log.WithError(err).Error("write response failed")

		return
	}

	logs.Log.
		WithField("user_id", tx.UserID).
		WithField("amount", tx.Amount).
		WithField("currency", tx.Currency).
		WithField("transaction_type", tx.TransactionType).
		WithField("category", tx.Category).
		WithField("transaction_date", tx.TransactionDate).
		WithField("description", tx.Description).
		WithField("commission", tx.Commission).
		Info("transaction successfully added")
}

func (h *Handlers) getTransactions(w http.ResponseWriter, r *http.Request) {
	userID, err := h.parseTokenFromRequest(r)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "parse token failed")

		return
	}

	if userID == "" {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "user id is required")

		return
	}

	transactions := make([]*models.Transactions, 0)
	result := h.gorm.Where("user_id = ?", userID).Find(&transactions)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "get translations failed")

		return
	}

	logEntry := logs.Log.WithField("user_id", userID)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		logEntry.WithError(err).Error("write response failed")

		return
	}

	logEntry.Info("get transactions successfully")
}

func (h *Handlers) getTransactionsWithCurrency(w http.ResponseWriter, r *http.Request, requestedCurrency string) {
	userID, err := h.parseTokenFromRequest(r)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "parse token failed")

		return
	}

	if userID == "" {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "user id is required")

		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "parse id failed")

		return
	}

	if id == 0 {
		logs.HandleMessage(w, r, http.StatusBadRequest, "id is required")

		return
	}

	var originalTx models.Transactions
	result := h.gorm.Where("user_id = ? AND id = ?", userID, id).First(&originalTx)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "find failed")

		return
	}

	converters := make(map[string]string)
	converters["base_currency"] = originalTx.Currency
	converters["currencies"] = requestedCurrency

	rates, err := currencies.LatestData(converters)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "get rates failed")

		return
	}

	if rates == nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "rates is nil")

		return
	}

	logEntry := logs.Log.
		WithField("id", id).
		WithField("user_id", userID).
		WithField("base_currency", originalTx.Currency).
		WithField("requested_currency", requestedCurrency)

	rate := rates.Data[requestedCurrency]

	tx := &models.ConvertedTransaction{
		Amount:       originalTx.Amount * rate,
		Transactions: originalTx,
		Currency:     requestedCurrency,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logEntry.WithError(err).Error("write response failed")

		return
	}

	logEntry.Info("get transactions with currency successfully")
}

func (h *Handlers) updateTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := h.parseTokenFromRequest(r)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "parse token failed")

		return
	}

	if userID == "" {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "user id is required")

		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "parse id failed")

		return
	}

	if id == 0 {
		logs.HandleMessage(w, r, http.StatusBadRequest, "id is required")

		return
	}

	var tx models.Transactions
	result := h.gorm.Where("user_id = ? AND id = ?", userID, id).First(&tx)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusNotFound, "transaction not found")

		return
	}

	err = json.NewDecoder(r.Body).Decode(&tx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.calculateCommission(&tx)
	tx.TransactionDate = time.Now()

	result = h.gorm.Save(&tx)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "update tx failed")

		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(tx)
	if err != nil {
		logs.Log.WithError(err).Error("write response failed")

		return
	}

	logs.Log.
		WithField("user_id", tx.UserID).
		WithField("amount", tx.Amount).
		WithField("currency", tx.Currency).
		WithField("transaction_type", tx.TransactionType).
		WithField("category", tx.Category).
		WithField("transaction_date", tx.TransactionDate).
		WithField("description", tx.Description).
		WithField("commission", tx.Commission).
		Info("transaction successfully updated")
}

func (h *Handlers) deleteTransaction(w http.ResponseWriter, r *http.Request) {
	userID, err := h.parseTokenFromRequest(r)
	if err != nil {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "parse token failed")

		return
	}

	if userID == "" {
		logs.HandleMessage(w, r, http.StatusUnauthorized, "user id is required")

		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		logs.HandleMessage(w, r, http.StatusBadRequest, "parse id failed")

		return
	}

	if id == 0 {
		logs.HandleMessage(w, r, http.StatusBadRequest, "id is required")

		return
	}

	result := h.gorm.Where("user_id = ?", userID).Delete(&models.Transactions{}, id)
	if result.Error != nil {
		logs.HandleMessage(w, r, http.StatusInternalServerError, "delete transaction failed")

		return
	}

	if result.RowsAffected == 0 {
		logs.HandleMessage(w, r, http.StatusOK, "records not found")

		return
	}

	w.WriteHeader(http.StatusOK)

	logs.Log.
		WithField("id", id).
		WithField("user_id", userID).
		Info("transaction successfully deleted")
}

func (h *Handlers) calculateCommission(c *models.Transactions) {
	switch {
	case c.TransactionType == "transfer" && c.Currency == "USD", c.Currency == "EUR", c.Currency == "GBP", c.Currency == "JPY":
		c.Commission = c.Amount * 0.02
	case c.TransactionType == "transfer" && c.Currency == "RUB":
		c.Commission = c.Amount * 0.05
	case c.TransactionType == "purchase", c.TransactionType == "top-up":
		c.Commission = 0
	default:
		c.Commission = 0
	}

}

func (h *Handlers) parseTokenFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header is required")

	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	jwtToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("YourSigningKey"), nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID := fmt.Sprintf("%.0f", claims["user_id"].(float64))
	return userID, nil
}

func New(gorm *gorm.DB) *Handlers {
	return &Handlers{
		gorm: gorm,
	}
}
