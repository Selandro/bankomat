package model

import (
	"errors"
	"sync"
)

type BankAccount interface {
	Deposit(amount float64) error
	Withdraw(amount float64) error
	GetBalance() float64
}

type Account struct {
	ID      int
	Balance float64
	mu      sync.RWMutex // Мьютекс для обеспечения потокобезопасности
}

func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return errors.New("сумма депозита должна быть положительной")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Balance = a.Balance + amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("сумма снятия должна быть положительной")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Balance < amount {
		return errors.New("недостаточно средств")
	}
	a.Balance -= amount
	return nil
}

func (a *Account) GetBalance() float64 {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.Balance
}
