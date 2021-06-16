package rabbitmq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/models"
)

func SendNFTToRabbit(nft models.NFT) {
	amqpChannel := getRabbitChannel()

	defer func() {
		closeErr := amqpChannel.Close()

		if closeErr != nil {
			fmt.Println("")
		}
	}()

	body, err := json.Marshal(nft)

	if err != nil {
		fmt.Println(fmt.Sprintf("error encoding JSON: %s", err))
		return
	}

	queue, err := amqpChannel.QueueDeclare(
		config.QueueIndexing,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(fmt.Sprintf("error during declare `QueueIndexing` queue: %s", err))
		return
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		fmt.Println(fmt.Sprintf("error publishing message: %s", err))
		return
	}

	if !nft.IsFinite {
		fmt.Println(fmt.Sprintf("Sent nft to Rabbit: id = %s, addr = %s ", nft.NFTID, nft.ContractAddress))
	}
}
