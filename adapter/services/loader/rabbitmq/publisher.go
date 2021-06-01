package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/models"
	"log"
)

func SendNFTToRabbit(nft models.NFT) {
	conn := getRabbitConn()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	body, err := json.Marshal(nft)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}
	queue, err := amqpChannel.QueueDeclare(config.QueueIndexing, true, false, false, false, nil)
	handleError(err, "Could not declare `QueueIndexing` queue")

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain", //todo not sure of this content type
		Body:         body,
	})

	if err != nil {
		log.Fatalf("Error publishing message: %s", err)
	}

	log.Printf("Sent nft to Rabbit: id = %s, addr = %s ", nft.NFTID, nft.ContractAddress)
}
