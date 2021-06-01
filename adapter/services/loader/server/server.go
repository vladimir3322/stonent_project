package server

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/events"
	"github.com/vladimir3322/stonent_go/tools/erc1155"
	"log"
	"math/big"
	"net/http"
	"strconv"
)

func getImageSource(w http.ResponseWriter, r *http.Request) {
	address := r.URL.Query().Get("address")

	if address == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Address query is required")
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Id query is required")
		return
	}

	imageId := new(big.Int)
	imageId, succeedParsedId := imageId.SetString(id, 10)

	if !succeedParsedId {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Invalid id")
		return
	}

	conn, ethErr := ethclient.Dial(config.ProviderUrl)

	if ethErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprint(ethErr))
		return
	}

	contract, ercErr := erc1155.NewErc1155(common.HexToAddress(address), conn)

	if ercErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprint(ercErr))
		return
	}

	imageSource, downErr := events.GetById(contract, imageId)

	if downErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprint(downErr))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(imageSource))
}

func getStatistics(w http.ResponseWriter, _ *http.Request) {
	type IResponse struct {
		CountOfFound int
		CountOfDownloaded int
	}

	events.Mutex.Lock()
	defer events.Mutex.Unlock()

	var response = IResponse{
		CountOfFound: events.CountOfFound,
		CountOfDownloaded: events.CountOfDownloaded,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func Run() {
	http.HandleFunc("/image_source", getImageSource)
	http.HandleFunc("/statistics", getStatistics)

	err := http.ListenAndServe(":" + strconv.Itoa(config.ServerPort), nil)

	if err != nil {
		log.Fatal("Server starting:", err)
		return
	}

	fmt.Println("Loader server is here!")
}
