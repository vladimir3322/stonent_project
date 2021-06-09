package events

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladimir3322/stonent_go/ml"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/tools/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
)

func getImageSource(ipfsHost string, ipfsPath string) (string, error) {
	ipfsMetadataUrl := ipfsHost + ipfsPath
	imageMetadataRes, err := http.Get(ipfsMetadataUrl)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsMetadataUrl, err))
	}
	if imageMetadataRes.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("error with %s invalid response code: %d", ipfsMetadataUrl, imageMetadataRes.StatusCode))
	}

	defer func() {
		err := imageMetadataRes.Body.Close()

		if err != nil {
			fmt.Println(err)
		}
	}()

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.NewDecoder(imageMetadataRes.Body).Decode(&jsonBody)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if imageMetadataParserErr != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsMetadataUrl, imageMetadataParserErr))
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", ipfsMetadataUrl, err))
	}

	imageSourceUrl := ipfsHost + "/ipfs" + parsedImageUrl.Path
	imageSourceRes, err := http.Get(imageSourceUrl)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s %s: %s", ipfsMetadataUrl, imageSourceUrl, err))
	}
	if imageSourceRes.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("error with %s %s: invalid response code %d", ipfsMetadataUrl, imageSourceUrl, imageSourceRes.StatusCode))
	}

	defer func() {
		err := imageSourceRes.Body.Close()

		if err != nil {
			fmt.Println(err)
		}
	}()

	imageSource, err := ioutil.ReadAll(imageSourceRes.Body)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New(fmt.Sprintf("error with %s: %s", imageSourceUrl, err))
	}

	b64ImageSource := base64.StdEncoding.EncodeToString(imageSource)

	return b64ImageSource, nil
}

func downloadImageWithWaiter(address string, nftId string, ipfsHost string, ipfsPath string, waiter *sync.WaitGroup, cb func(isSucceed bool)) {
	defer waiter.Done()

	isSucceed := downloadImage(address, nftId, ipfsHost, ipfsPath)
	cb(isSucceed)
}

func downloadImage(address string, nftId string, ipfsHost string, ipfsPath string) bool {
	imageSource, err := getImageSource(ipfsHost, ipfsPath)

	if err != nil {
		ml.SentRejectedImageByIPFS(address, nftId, ipfsHost, ipfsPath, err)

		return false
	}

	rabbitmq.SendNFTToRabbit(models.NFT{
		NFTID: nftId,
		ContractAddress: address,
		Data: imageSource,
		IsFinite: false,
	})

	return true
}
