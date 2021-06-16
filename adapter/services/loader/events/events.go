package events

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/eth"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"math/big"
	"strings"
	"sync"
)

type iImageMetadata struct {
	Image string
}

func GetById(contract *erc1155.Erc1155, id *big.Int) (string, error) {
	opt := &bind.FilterOpts{}
	s := []*big.Int{id}
	event, err := contract.FilterURI(opt, s)

	if err != nil {
		return "", err
	}

	isExist := event.Next()

	if !isExist {
		return "", errors.New("event not found")
	}

	ipfsPath := strings.ReplaceAll(event.Event.Value, "/ipfs/", "")
	imageSource, err := getImageSource(ipfsPath)

	if err != nil {
		return "", err
	}

	return imageSource, nil
}

func GetEvents(address string, contract *erc1155.Erc1155, startBlock uint64, endBlock uint64, waiter *sync.WaitGroup) {
	defer waiter.Done()

	if IsExceededImagesLimitCount() {
		waiter.Done()
		return
	}

	if startBlock <= endBlock {
		var s []*big.Int

		opt := &bind.FilterOpts{Start: startBlock, End: &endBlock}
		past, err := contract.FilterURI(opt, s)

		if IsExceededImagesLimitCount() {
			waiter.Done()
			return
		}

		if err != nil {
			var middleBlock = (startBlock + endBlock) / 2

			waiter.Add(1)
			go GetEvents(address, contract, startBlock, middleBlock, waiter)
			waiter.Add(1)
			go GetEvents(address, contract, middleBlock+1, endBlock, waiter)
			return
		}

		notEmpty := true

		for notEmpty {
			if IsExceededImagesLimitCount() {
				waiter.Done()
				return
			}

			notEmpty = past.Next()

			if notEmpty {
				waiter.Add(1)

				ipfsPath := strings.ReplaceAll(past.Event.Value, "/ipfs/", "")

				go pushToBuffer(BufferItem{
					address:  address,
					nftId:    past.Event.Id.String(),
					ipfsPath: ipfsPath,
					waiter:   waiter,
				})
			}
		}
	} else {
		return
	}
}

func listenByWatcher(ethClient *ethclient.Client, address string, startBlock uint64) {
	contract, err := erc1155.NewErc1155(common.HexToAddress(address), ethClient)

	if err != nil {
		fmt.Println(fmt.Sprintf("failed listening events: %s", err))
		return
	}

	var s []*big.Int

	opts := &bind.WatchOpts{Start: &startBlock}
	ch := make(chan *erc1155.Erc1155URI)
	watcher, err := contract.WatchURI(opts, ch, s)

	if err != nil {
		fmt.Println(fmt.Sprintf("failed listening events: %s", err))
		return
	}

	for {
		select {
		case err := <-watcher.Err():
			fmt.Println(fmt.Sprintf("failed listening events: %s, restarting", err))
			return
		case Event := <-ch:
			fmt.Println(fmt.Sprintf("received event from listening: %s", Event.Id))

			ipfsPath := strings.ReplaceAll(Event.Value, "/ipfs/", "")

			go downloadImage(address, Event.Id.String(), ipfsPath)
		}
	}
}

func ListenEvents(address string, startBlock uint64) {
	for {
		ethClient := eth.GetEthClient()

		listenByWatcher(ethClient, address, startBlock)

		startBlock = eth.GetLatestBlockNumber(ethClient)
	}
}
