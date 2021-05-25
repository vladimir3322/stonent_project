package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/rabbitmq"
	"github.com/vladimir3322/stonent_go/tools/models"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
)

func getImageSource(ipfsHost string, ipfsPath string) ([]byte, error) {
	ipfsMetadataUrl := ipfsHost + ipfsPath
	imageMetadataRes, err := http.Get(ipfsMetadataUrl)

	if config.DownloadImageMaxCount != -1 && countOfDownloaded >= config.DownloadImageMaxCount {
		return nil, errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return nil, errors.New("error with " + ipfsMetadataUrl + " : " + fmt.Sprint(err))
	}
	if imageMetadataRes.StatusCode != http.StatusOK {
		return nil, errors.New("error with " + ipfsMetadataUrl + " invalid response code: " + strconv.Itoa(imageMetadataRes.StatusCode))
	}

	defer imageMetadataRes.Body.Close()

	var jsonBody iImageMetadata
	imageMetadataParserErr := json.NewDecoder(imageMetadataRes.Body).Decode(&jsonBody)

	if config.DownloadImageMaxCount != -1 && countOfDownloaded >= config.DownloadImageMaxCount {
		return nil, errors.New("downloaded image limit exceeded")
	}
	if imageMetadataParserErr != nil {
		return nil, errors.New("error with: " + ipfsMetadataUrl + " : " + fmt.Sprint(imageMetadataParserErr))
	}

	parsedImageUrl, err := url.Parse(jsonBody.Image)

	if config.DownloadImageMaxCount != -1 && countOfDownloaded >= config.DownloadImageMaxCount {
		return nil, errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return nil, errors.New("error with: " + ipfsMetadataUrl + fmt.Sprint(err))
	}

	imageSourceUrl := ipfsHost + "/ipfs" + parsedImageUrl.Path
	imageSourceRes, err := http.Get(imageSourceUrl)

	if config.DownloadImageMaxCount != -1 && countOfDownloaded >= config.DownloadImageMaxCount {
		return nil, errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return nil, errors.New("error with: " + ipfsMetadataUrl + " " + imageSourceUrl + " : " + fmt.Sprint(err))
	}
	if imageSourceRes.StatusCode != http.StatusOK {
		return nil, errors.New("Error with: " + ipfsMetadataUrl + " " + imageSourceUrl + " invalid response code: " + strconv.Itoa(imageSourceRes.StatusCode))
	}

	defer imageSourceRes.Body.Close()

	imageSource, err := ioutil.ReadAll(imageSourceRes.Body)

	if config.DownloadImageMaxCount != -1 && countOfDownloaded >= config.DownloadImageMaxCount {
		return nil, errors.New("downloaded image limit exceeded")
	}
	if err != nil {
		return nil, errors.New("error with: " + imageSourceUrl + " " + fmt.Sprint(err))
	}

	return imageSource, nil
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

	img := string(imageSource[:])
	fmt.Println("Downloaded image size:", len(img))
	fmt.Println("Image fragment:", img[len(img)-100:])

	file, err := os.Create(nftId)

	if err != nil {
		fmt.Println("ERROR WITH SAVING IMAGE:", err)
	}

	defer file.Close()
	_, saveErr := file.WriteString(string(imageSource[:]))

	if saveErr != nil {
		fmt.Println("ERROR WITH SAVING IMAGE:", err)
	}

	rabbitmq.SendNFTToRabbit(models.NFT{
		NFTID:           nftId,
		ContractAddress: address,
		Data:            string(imageSource[:]),
	})

	return true
}
