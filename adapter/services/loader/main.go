package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/eth"
	"github.com/vladimir3322/stonent_go/events"
	"github.com/vladimir3322/stonent_go/ipfs"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/server"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"github.com/vladimir3322/stonent_go/tools/models"
	"github.com/vladimir3322/stonent_go/tools/utils"
	"sync"
)

func main() {
	go server.Run()
	rabbitmq.InitRabbit()
	ipfs.Init()

	var contractsAddresses = []string{"0xd07dc4262bcdbf85190c01c996b4c06a461d2430"}
	var completedContracts = 0

	ethConnection := eth.GetEthClient()

	// Все картины
	//startBlockNumber := 0
	//latestBlockNumber := eth.GetLatestBlockNumber(ethConnection)

	// Много картин
	startBlockNumber := uint64(12291943)
	latestBlockNumber := eth.GetLatestBlockNumber(ethConnection)

	// 4 картины
	//startBlockNumber := 12291940
	//latestBlockNumber := 12291943

	for _, contractAddress := range contractsAddresses {
		go getEvents(ethConnection, contractAddress, startBlockNumber, latestBlockNumber, func() {
			completedContracts += 1

			if len(contractsAddresses) == completedContracts {
				fmt.Println("loader completed successfully")

				rabbitmq.SendNFTToRabbit(models.NFT{
					IsFinite: true,
				})
			}
		})

		go events.ListenEvents(contractAddress, latestBlockNumber)
	}

	utils.WaitSignals()
}

func getEvents(ethConnection *ethclient.Client, address string, startBlock uint64, endBlock uint64, onFinished func()) {
	contract, err := erc1155.NewErc1155(common.HexToAddress(address), ethConnection)

	if err != nil {
		fmt.Println(fmt.Sprintf("whoops something went wrong: %s", err))
	}

	var waiter = &sync.WaitGroup{}

	waiter.Add(1)
	events.GetEvents(address, contract, startBlock, endBlock, waiter)
	go events.RunBuffer()
	waiter.Wait()
	onFinished()
}
