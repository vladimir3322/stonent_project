package models

type NFT struct {
	NFTID string `json:"nftID"`
	ContractAddress string `json:"contractAddress"`
	Data string `json:"data"`
	IsFinite bool `json:"isFinite"`
}
