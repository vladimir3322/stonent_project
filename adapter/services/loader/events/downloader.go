package events

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladimir3322/stonent_go/ipfs"
	"github.com/vladimir3322/stonent_go/ml"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/tools/models"
	"net/url"
	"sync"
)

func getImageSource(ipfsPath string) (string, error) {
	imageMetadataRes, err := ipfs.Get(ipfsPath)

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsPath, err))
	}

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.Unmarshal(imageMetadataRes, &jsonBody)

	if imageMetadataParserErr != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsPath, imageMetadataParserErr))
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsPath, err))
	}

	imageSourceRes, err := ipfs.Get(parsedImageUrl.Path[1:])

	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s %s: %s", ipfsPath, parsedImageUrl.Path, err))
	}

	b64ImageSource := base64.StdEncoding.EncodeToString(imageSourceRes)

	return b64ImageSource, nil
}

func downloadImageWithWaiter(address string, nftId string, ipfsPath string, waiter *sync.WaitGroup, cb func(isSucceed bool)) {
	defer waiter.Done()

	isSucceed := downloadImage(address, nftId, ipfsPath)

	if isSucceed {
		CountOfDownloaded += 1
	}

	cb(isSucceed)
}

func downloadImage(address string, nftId string, ipfsPath string) bool {
	imageSource, err := getImageSource(ipfsPath)

	if err != nil {
		ml.SentRejectedImageByIPFS(address, nftId, ipfsPath, err)

		return false
	}

	rabbitmq.SendNFTToRabbit(models.NFT{
		NFTID:           nftId,
		ContractAddress: address,
		Data:            imageSource,
		IsFinite:        false,
	})

	return true
}
