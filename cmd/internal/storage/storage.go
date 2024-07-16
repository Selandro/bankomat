package storage

import (
	"errors"
	"sync"

	"main.go/model"
)

var (
	AccountMap      map[int]*model.Account
	AccountMapMutex sync.RWMutex
)

// InitCache инициализирует карту аккаунтов.
func InitStorage() {
	AccountMap = make(map[int]*model.Account)

}

// AddAccount добавляет аккаунт в кэш.
func AddAccount(account *model.Account) error {
	AccountMapMutex.Lock()
	defer AccountMapMutex.Unlock()
	_, err := AccountMap[account.ID]
	if !err {
		AccountMap[account.ID] = account
		return nil
	} else {
		err := errors.New("аккаунт уже существует")
		return err

	}
}

// GetAccount возвращает аккаунт из кэша по ID.
func GetAccount(id int) (*model.Account, bool) {
	AccountMapMutex.RLock()
	defer AccountMapMutex.RUnlock()
	account, exists := AccountMap[id]
	return account, exists
}
