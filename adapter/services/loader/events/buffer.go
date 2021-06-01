package events

import (
	"github.com/vladimir3322/stonent_go/config"
	"sync"
)

type BufferItem struct {
	ipfsHost string
	ipfsPath string
	address  string
	nftId    string
	waiter   *sync.WaitGroup
}

var Buffer = make(chan BufferItem)
var bufferSize = 0
var CountOfFound = 0
var CountOfDownloaded = 0
var Mutex = sync.Mutex{}

func IsExceededImagesLimitCount() bool {
	return config.DownloadImageMaxCount != -1 && CountOfDownloaded >= config.DownloadImageMaxCount
}

func pushToBuffer(item BufferItem) {
	Mutex.Lock()
	defer Mutex.Unlock()

	CountOfFound += 1

	for bufferSize >= config.DownloadImageBufferSize {

	}

	if IsExceededImagesLimitCount() {
		item.waiter.Done()
		return
	}

	bufferSize += 1
	Buffer <- item
}

func RunBuffer() {
	for {
		select {
		case item := <-Buffer:
			if IsExceededImagesLimitCount() {
				item.waiter.Done()
				return
			}

			go downloadImageWithWaiter(item.address, item.nftId, item.ipfsHost, item.ipfsPath, item.waiter, func(isSucceed bool) {
				bufferSize -= 1

				if isSucceed {
					Mutex.Lock()
					defer Mutex.Unlock()

					CountOfDownloaded += 1
				}
			})
		}
	}
}
