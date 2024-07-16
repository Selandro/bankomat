package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"main.go/cmd/internal/storage"
	"main.go/model"
)

// Response представляет структуру ответа для операций
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// CreateAccountHandler обрабатывает создание нового аккаунта
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID int `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account := &model.Account{ID: req.ID, Balance: 0}

	err := storage.AddAccount(account)
	if err != nil {
		http.Error(w, "аккаунт уже существует", http.StatusBadRequest)
		return
	}

	loggerHandlers("CreateAccount", account.ID, 0)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

// DepositHandler обрабатывает пополнение баланса аккаунта
func DepositHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/accounts/") : len(r.URL.Path)-len("/deposit")]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, exists := storage.GetAccount(id)
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Канал для передачи ответа
	responseChan := make(chan Response)

	// Запускаем горутину для выполнения операции пополнения
	go func() {
		err := account.Deposit(req.Amount)
		if err != nil {
			// Передаем ошибку в канал
			responseChan <- Response{http.StatusBadRequest, err.Error()}
			return
		}

		loggerHandlers("Deposit", account.ID, req.Amount)
		// Передаем успешный ответ в канал
		responseChan <- Response{http.StatusOK, "Пополнение успешно"}
	}()

	// Получаем ответ из канала
	response := <-responseChan
	w.WriteHeader(response.Status)
	w.Write([]byte(response.Message))
}

// WithdrawHandler обрабатывает снятие средств с аккаунта
func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/accounts/") : len(r.URL.Path)-len("/withdraw")]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, exists := storage.GetAccount(id)
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Канал для передачи ответа
	responseChan := make(chan Response)

	// Запускаем горутину для выполнения операции снятия средств
	go func() {
		err := account.Withdraw(req.Amount)
		if err != nil {
			// Передаем ошибку в канал
			responseChan <- Response{http.StatusBadRequest, err.Error()}
			return
		}

		loggerHandlers("Withdraw", account.ID, req.Amount)
		// Передаем успешный ответ в канал
		responseChan <- Response{http.StatusOK, "Снятие успешно"}
	}()

	// Получаем ответ из канала
	response := <-responseChan
	w.WriteHeader(response.Status)
	w.Write([]byte(response.Message))
}

// GetBalanceHandler обрабатывает получение баланса аккаунта
func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/accounts/") : len(r.URL.Path)-len("/balance")]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid account ID", http.StatusBadRequest)
		return
	}

	// Получаем аккаунт из хранилища по ID
	account, exists := storage.GetAccount(id)
	if !exists {
		http.Error(w, "Account not found", http.StatusNotFound)
		return
	}

	// Канал для передачи баланса
	responseChan := make(chan float64)

	// Запускаем горутину для получения баланса
	go func() {
		balance := account.GetBalance()
		// Логируем операцию получения баланса
		loggerHandlers("GetBalance", account.ID, balance)
		// Передаем баланс в канал
		responseChan <- balance
	}()

	// Получаем баланс из канала
	balance := <-responseChan
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

// loggerHandlers логирует операции с аккаунтом
func loggerHandlers(operation string, accountID int, amount float64) {
	fmt.Printf("%s: Operation: %s, Account ID: %d, Amount: %.2f, Time: %s\n",
		time.Now().Format(time.RFC3339), operation, accountID, amount, time.Now().Format(time.RFC3339))
}
