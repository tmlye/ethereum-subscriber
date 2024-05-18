package storage

import (
	"sync"

	"github.com/tmlye/ethereum-subscriber/pkg/types"
)

type Store interface {
	Subscribe(address string) bool
	IsSubscribed(address string) bool
	AddTransaction(address string, tx types.Transaction)
	GetTransactions(address string) []types.Transaction
}

type MemoryStore struct {
	subscriptions map[string]bool
	transactions  map[string][]types.Transaction
	mutex         sync.Mutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		subscriptions: make(map[string]bool),
		transactions:  make(map[string][]types.Transaction),
	}
}

func (s *MemoryStore) Subscribe(address string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.subscriptions[address]; exists {
		return false
	}
	s.subscriptions[address] = true
	return true
}

func (s *MemoryStore) IsSubscribed(address string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.subscriptions[address]
}

func (s *MemoryStore) AddTransaction(address string, tx types.Transaction) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.transactions[address] = append(s.transactions[address], tx)
}

func (s *MemoryStore) GetTransactions(address string) []types.Transaction {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.transactions[address]
}
