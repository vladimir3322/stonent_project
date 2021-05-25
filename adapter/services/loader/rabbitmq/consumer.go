package rabbitmq

import (
	"encoding/json"
	"github.com/vladimir3322/stonent_go/config"
	"github.com/vladimir3322/stonent_go/tools/models"
	"github.com/vladimir3322/stonent_go/tools/utils"
	"log"
)

// for test purposes
func ConsumeEvents() {

	conn := InitRabbit()

	ch, err := conn.Channel()
	handleError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		config.QueueIndexing, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	handleError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	handleError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {

			var nftReceived models.NFT
			err := json.Unmarshal(d.Body, &nftReceived)
			if err != nil {
				handleError(err, "Error while unmarshalling received nft data")
			}

			log.Printf("Received a nft with id: %s", nftReceived.NFTID)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	utils.WaitSignals()
}
