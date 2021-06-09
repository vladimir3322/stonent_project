package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/events"
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

	var contractsAddresses = []string{"0xd07dc4262bcdbf85190c01c996b4c06a461d2430"}
	var completedContracts = 0

	for _, contractAddress := range contractsAddresses {
		//go getEvents(contractAddress, 12291940, 12291943, func () {
		//	completedContracts += 1
		//
		//	if len(contractsAddresses) == completedContracts {
		//		fmt.Println("loader completed successfully")
		//
		//		rabbitmq.SendNFTToRabbit(models.NFT{
		//			IsFinite: true,
		//		})
		//	}
		//}) // 4 картины
		go getEvents(contractAddress, 12291940, 12594330, func () {
			completedContracts += 1

			if len(contractsAddresses) == completedContracts {
				fmt.Println("loader completed successfully")

				rabbitmq.SendNFTToRabbit(models.NFT{
					IsFinite: true,
				})
			}
		}) // Много картин

		//go getEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 0, 12291943) // Почти все
		//go getEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 12211843, 12291943) // Много картин

		//go listenEvents("0xd07dc4262bcdbf85190c01c996b4c06a461d2430", 12291943)
	}

	utils.WaitSignals()

	//fmt.Println(api.GetLatestBlock(conn))
}

func getEvents(address string, startBlock uint64, endBlock uint64, onFinished func()) {
	conn, err := ethclient.Dial(config.ProviderUrl)

	if err != nil {
		fmt.Printf("whoops something went wrong: %s", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)

	if err != nil {
		fmt.Printf("whoops something went wrong: %s", err)
	}

	var waiter = &sync.WaitGroup{}

	waiter.Add(1)
	events.GetEvents(address, contract, startBlock, endBlock, waiter)
	go events.RunBuffer()
	waiter.Wait()
	onFinished()
}

func listenEvents(address string, startBlock uint64) {
	conn, err := ethclient.Dial(config.ProviderUrl)

	if err != nil {
		fmt.Printf("whoops something went wrong: %s", err)
	}

	contract, err := erc1155.NewErc1155(common.HexToAddress(address), conn)

	if err != nil {
		fmt.Printf("whoops something went wrong: %s", err)
	}

	events.ListenEvents(address, contract, startBlock)
}
