package blockpoller

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tmlye/ethereum-subscriber/pkg/ethgateway"
	"github.com/tmlye/ethereum-subscriber/pkg/storage"
)

type BlockPoller struct {
	gateway                  ethgateway.Gateway
	store                    storage.Store
	lastProcessedBlockNumber string
}

func NewBlockPoller(gateway ethgateway.Gateway, store storage.Store) *BlockPoller {
	return &BlockPoller{
		gateway:                  gateway,
		store:                    store,
		lastProcessedBlockNumber: "0x0",
	}
}

func (p *BlockPoller) LastProcessedBlock() int64 {
	return hexToInt(p.lastProcessedBlockNumber)
}

func (p *BlockPoller) ProcessBlock(blockNumber string) error {
	log.Println("Processing block", hexToInt(blockNumber))
	block, err := p.gateway.GetBlockByNumber(blockNumber)
	if err != nil {
		log.Println("Could not fetch block", blockNumber)
		return err
	}

	for _, tx := range block.Transactions {
		if p.store.IsSubscribed(tx.From) {
			log.Println("Adding tx", tx)
			p.store.AddTransaction(tx.From, tx)
		}
		if p.store.IsSubscribed(tx.To) {
			log.Println("Adding tx", tx)
			p.store.AddTransaction(tx.To, tx)
		}
	}

	return nil
}

func (p *BlockPoller) PollBlocks(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		currentBlock, err := p.gateway.GetCurrentBlock()
		if err != nil {
			log.Println("Error getting current block, continuing")
			continue
		}

		if currentBlock == p.lastProcessedBlockNumber {
			log.Println("Already processed block, continuing", currentBlock)
			continue
		}

		err = p.ProcessBlock(currentBlock)
		if err != nil {
			log.Println("Error parsing block, continuing", currentBlock)
			continue
		}

		p.lastProcessedBlockNumber = currentBlock
	}
}

func hexToInt(hexNum string) int64 {
	if strings.HasPrefix(hexNum, "0x") {
		hexNum = hexNum[2:]
	}

	num, err := strconv.ParseInt(hexNum, 16, 64)
	if err != nil {
		log.Println("Error:", err)
		return 0
	}

	return num
}
