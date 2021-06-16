package ml

import (
	"fmt"
	"github.com/vladimir3322/stonent_go/config"
	"net/http"
	"net/url"
)

func SentRejectedImageByIPFS(contractAddress string, nftId string, ipfsPath string, error error) {
	query := fmt.Sprintf(
		"?contract_address=%s&nft_id=%s&ipfs_path=%s&error=%s",
		contractAddress,
		nftId,
		url.QueryEscape(ipfsPath),
		url.QueryEscape(fmt.Sprint(error)),
	)
	requestUrl := fmt.Sprintf("%s%s%s", config.MlUrl, "/register_rejected_image_by_ipfs", query)

	res, err := http.Get(requestUrl)

	if err != nil {
		fmt.Println(fmt.Sprintf("fail to request to ML service: %s", err))
		return
	}
	if res.StatusCode != http.StatusOK {
		fmt.Println(fmt.Sprintf("fail to request to ML service: %d", res.StatusCode))
		return
	}
}
