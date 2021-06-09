package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
	"github.com/vladimir3322/stonent_go/config"
	"time"
)

var rabbitConn *amqp.Connection

func InitRabbit() *amqp.Connection {
	var err error
	rabbitAddr := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.RabbitLogin, config.RabbitPass, config.RabbitHost, config.RabbitPort)
	rabbitConn, err = amqp.Dial(rabbitAddr)

	if err != nil {
		fmt.Println("Rabbit connection failed, restarting")
		time.Sleep(time.Second)
		return InitRabbit()
	}

	fmt.Println("Successfully connected to Rabbit")
	return rabbitConn
}

func getRabbitConn() *amqp.Connection {
	if rabbitConn == nil {
		return InitRabbit()
	}
	return rabbitConn
}
