package events

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"math/big"
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

	// TODO: set to 0
	imageSource, err := getImageSource(config.IpfsLink[len(config.IpfsLink)-1], event.Event.Value)

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
		opt := &bind.FilterOpts{Start: startBlock, End: &endBlock}
		s := []*big.Int{}
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
		ipfsNodeIndex := 0

		for notEmpty {
			if IsExceededImagesLimitCount() {
				waiter.Done()
				return
			}

			notEmpty = past.Next()
			if notEmpty {
				waiter.Add(1)

				go pushToBuffer(BufferItem{
					address:  address,
					nftId:    past.Event.Id.String(),
					ipfsHost: config.IpfsLink[ipfsNodeIndex],
					ipfsPath: past.Event.Value,
					waiter:   waiter,
				})

				ipfsNodeIndex += 1
				ipfsNodeIndex %= len(config.IpfsLink)
			}
		}
	} else {
		return
	}
}

func ListenEvents(address string, contract *erc1155.Erc1155, startBlock uint64) {
	s := []*big.Int{}
	ch := make(chan *erc1155.Erc1155URI)
	opts := &bind.WatchOpts{Start: &startBlock}
	watcher, err := contract.WatchURI(opts, ch, s)

	if err != nil {
		fmt.Println("Failed listening events:", err)
	}

	ipfsNodeIndex := 0

	for {
		select {
		case err := <-watcher.Err():
			fmt.Println("Failed listening events:", err)
		case Event := <-ch:
			fmt.Println(Event.Value)

			go downloadImage(address, Event.Id.String(), config.IpfsLink[ipfsNodeIndex], Event.Value)

			ipfsNodeIndex += 1
			ipfsNodeIndex %= len(config.IpfsLink)
		}
	}
}
