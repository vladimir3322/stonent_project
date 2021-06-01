package events

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/tools/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
)

func getImageSource(ipfsHost string, ipfsPath string) (string, error) {
	ipfsMetadataUrl := ipfsHost + ipfsPath
	imageMetadataRes, err := http.Get(ipfsMetadataUrl)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New("error with " + ipfsMetadataUrl + " : " + fmt.Sprint(err))
	}
	if imageMetadataRes.StatusCode != http.StatusOK {
		return "", errors.New("error with " + ipfsMetadataUrl + " invalid response code: " + strconv.Itoa(imageMetadataRes.StatusCode))
	}

	defer imageMetadataRes.Body.Close()

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.NewDecoder(imageMetadataRes.Body).Decode(&jsonBody)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if imageMetadataParserErr != nil {
		return "", errors.New("error with: " + ipfsMetadataUrl + " : " + fmt.Sprint(imageMetadataParserErr))
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New("error with: " + ipfsMetadataUrl + fmt.Sprint(err))
	}

	imageSourceUrl := ipfsHost + "/ipfs" + parsedImageUrl.Path
	imageSourceRes, err := http.Get(imageSourceUrl)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New("error with: " + ipfsMetadataUrl + " " + imageSourceUrl + " : " + fmt.Sprint(err))
	}
	if imageSourceRes.StatusCode != http.StatusOK {
		return "", errors.New("Error with: " + ipfsMetadataUrl + " " + imageSourceUrl + " invalid response code: " + strconv.Itoa(imageSourceRes.StatusCode))
	}

	defer imageSourceRes.Body.Close()

	imageSource, err := ioutil.ReadAll(imageSourceRes.Body)

	if IsExceededImagesLimitCount() {
		return "", errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return "", errors.New("error with: " + imageSourceUrl + " " + fmt.Sprint(err))
	}

	b64ImageSource := base64.StdEncoding.EncodeToString(imageSource);

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
		fmt.Println(err)
		return false
	}

	rabbitmq.SendNFTToRabbit(models.NFT{
		NFTID:           nftId,
		ContractAddress: address,
		Data:            imageSource,
	})

	return true
}
