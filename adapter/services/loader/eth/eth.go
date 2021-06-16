package eth

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/api"
	"github.com/vladimir3322/stonent_go/tools/models"
	"time"
)

func GetEthClient() *ethclient.Client {
	ethClient, err := ethclient.Dial(config.ProviderUrl)

	if err != nil {
		time.Sleep(time.Millisecond * 100)
		return GetEthClient()
	}

	return ethClient
}

func getLatestBlock(ethClient *ethclient.Client) *models.Block {
	latestBlock := api.GetLatestBlock(ethClient)

	if latestBlock != nil {
		return latestBlock
	}

	time.Sleep(time.Millisecond * 100)
	return getLatestBlock(ethClient)
}

func GetLatestBlockNumber(ethClient *ethclient.Client) uint64 {
	latestBlock := getLatestBlock(ethClient)
	latestBlockNumber := uint64(latestBlock.BlockNumber)

	return latestBlockNumber
}
