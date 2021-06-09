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
		_, err := fmt.Fprintf(w, "address query is required")

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "id query is required")

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	imageId := new(big.Int)
	imageId, succeedParsedId := imageId.SetString(id, 10)

	if !succeedParsedId {
		w.WriteHeader(http.StatusBadRequest)
		_, err := fmt.Fprintf(w, "invalid id")

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	conn, ethErr := ethclient.Dial(config.ProviderUrl)

	if ethErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, fmt.Sprint(ethErr))

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	contract, ercErr := erc1155.NewErc1155(common.HexToAddress(address), conn)

	if ercErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, fmt.Sprint(ercErr))

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	imageSource, downErr := events.GetById(contract, imageId)

	if downErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err := fmt.Fprintf(w, fmt.Sprint(downErr))

		if err != nil {
			fmt.Println(err)
		}

		return
	}

	w.WriteHeader(http.StatusOK)

	_, err := w.Write([]byte(imageSource))

	if err != nil {
		fmt.Printf("error with server response writing: %s", err)
	}
}

func getStatistics(w http.ResponseWriter, _ *http.Request) {
	type IResponse struct {
		CountOfFound int
		CountOfDownloaded int
		CountOfRejected int
	}

	var response = IResponse{
		CountOfFound: events.CountOfFound,
		CountOfDownloaded: events.CountOfDownloaded,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		fmt.Println(err)
	}
}

func Run() {
	http.HandleFunc("/image_source", getImageSource)
	http.HandleFunc("/statistics", getStatistics)

	err := http.ListenAndServe(":" + strconv.Itoa(config.ServerPort), nil)

	if err != nil {
		log.Fatal("server starting:", err)
		return
	}

	fmt.Println("loader server is here!")
}
