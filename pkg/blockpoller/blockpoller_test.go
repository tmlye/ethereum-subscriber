package blockpoller

import (
	"testing"

	"github.com/tmlye/ethereum-subscriber/pkg/storage"
	"github.com/tmlye/ethereum-subscriber/pkg/types"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockGateway struct {
	mock.Mock
}

func (m *MockGateway) GetCurrentBlock() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *MockGateway) GetBlockByNumber(blockNumber string) (types.Block, error) {
	args := m.Called(blockNumber)
	return args.Get(0).(types.Block), args.Error(1)
}

func TestBlockPoller_ProcessBlock_ShouldStoreForTo(t *testing.T) {
	// Arrange
	currentBlock := "0x2"
	mockGateway := new(MockGateway)
	mockGateway.On("GetCurrentBlock").Return(currentBlock, nil).Once()
	transactions := []types.Transaction{
		{From: "0x456", To: "0x123", Value: "100"},
		{From: "0x456", To: "0xabc", Value: "200"},
	}
	block := types.Block{Transactions: transactions}
	mockGateway.On("GetBlockByNumber", currentBlock).Return(block, nil).Once()

	store := storage.NewMemoryStore()
	store.Subscribe("0x123")

	blockPoller := NewBlockPoller(mockGateway, store)

	// Act
	blockPoller.ProcessBlock(currentBlock)

	// Assert
	txs123 := store.GetTransactions("0x123")
	assert.Len(t, txs123, 1)
	assert.Equal(t, transactions[0], txs123[0])
	txs456 := store.GetTransactions("0x456")
	assert.Len(t, txs456, 0)
}

func TestBlockPoller_ProcessBlock_ShouldStoreForFrom(t *testing.T) {
	// Arrange
	currentBlock := "0x2"
	mockGateway := new(MockGateway)
	mockGateway.On("GetCurrentBlock").Return(currentBlock, nil).Once()
	transactions := []types.Transaction{
		{From: "0x123", To: "0x456", Value: "100"},
		{From: "0x456", To: "0xabc", Value: "200"},
	}
	block := types.Block{Transactions: transactions}
	mockGateway.On("GetBlockByNumber", currentBlock).Return(block, nil).Once()

	store := storage.NewMemoryStore()
	store.Subscribe("0x123")

	blockPoller := NewBlockPoller(mockGateway, store)

	// Act
	blockPoller.ProcessBlock(currentBlock)

	// Assert
	txs123 := store.GetTransactions("0x123")
	assert.Len(t, txs123, 1)
	assert.Equal(t, transactions[0], txs123[0])
	txs456 := store.GetTransactions("0x456")
	assert.Len(t, txs456, 0)
}

func TestBlockPoller_ProcessBlock_ShouldEmptyWhenNoTx(t *testing.T) {
	// Arrange
	currentBlock := "0x2"
	mockGateway := new(MockGateway)
	mockGateway.On("GetCurrentBlock").Return(currentBlock, nil).Once()
	block := types.Block{Transactions: []types.Transaction{}}
	mockGateway.On("GetBlockByNumber", currentBlock).Return(block, nil).Once()

	store := storage.NewMemoryStore()
	store.Subscribe("0x123")

	blockPoller := NewBlockPoller(mockGateway, store)

	// Act
	blockPoller.ProcessBlock(currentBlock)

	// Assert
	txs123 := store.GetTransactions("0x123")
	assert.Len(t, txs123, 0)
}

func TestBlockPoller_hexToInt(t *testing.T) {
	// Arrange
	tests := []struct {
		hexStr    string
		expected  int64
		shouldErr bool
	}{
		{"0x1", 1, false},
		{"0xA", 10, false},
		{"0x10", 16, false},
		{"0x64", 100, false},
		{"0xinvalid", 0, true},
	}

	for _, test := range tests {
		// Act
		result := hexToInt(test.hexStr)

		// Assert
		if test.shouldErr {
			assert.Equal(t, int64(0), result, "expected error for input %s", test.hexStr)
		} else {
			assert.Equal(t, test.expected, result, "expected %d for input %s, got %d", test.expected, test.hexStr, result)
		}
	}
}
