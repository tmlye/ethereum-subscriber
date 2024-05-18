package ethgateway

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/tmlye/ethereum-subscriber/pkg/types"
)

type Gateway interface {
	GetCurrentBlock() (string, error)
	GetBlockByNumber(string) (types.Block, error)
}

type EthGateway struct {
	url string
}

func NewEthGateway(url string) *EthGateway {
	return &EthGateway{
		url,
	}
}

func (gw *EthGateway) sendRPCRequest(method string, params []interface{}) ([]byte, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}
	requestBody, _ := json.Marshal(request)

	resp, err := http.Post(gw.url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (gw *EthGateway) GetCurrentBlock() (string, error) {
	response, err := gw.sendRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Println("Error fetching current block:", err)
		return "", err
	}

	var result struct {
		Result string `json:"result"`
	}
	if err := json.Unmarshal(response, &result); err != nil {
		log.Println("Error parsing response:", err)
		return "", err
	}

	return result.Result, nil
}

func (gw *EthGateway) GetBlockByNumber(blockNumber string) (types.Block, error) {
	blockData, err := gw.sendRPCRequest("eth_getBlockByNumber", []interface{}{blockNumber, true})
	if err != nil {
		log.Println("Error fetching block data:", err)
		return types.Block{}, err
	}

	var result struct {
		Block types.Block `json:"result"`
	}
	if err := json.Unmarshal(blockData, &result); err != nil {
		log.Println("Error parsing block data:", err)
		return types.Block{}, err
	}

	return result.Block, nil
}
