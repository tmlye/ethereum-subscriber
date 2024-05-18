package types

type Block struct {
	Transactions []Transaction `json:"transactions"`
}
