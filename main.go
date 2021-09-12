package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/kafka"
)

var (
	wg sync.WaitGroup
)

func main() {
	c, err := kafka.NewConsumer("plogger-kafka", "plogger-group", []string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatal(err)
	}

	consumer := kafka.Consumer{
		Messages: make(chan *sarama.ConsumerMessage, 1000),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.Consume(ctx, &consumer)
	for message := range consumer.Messages {

		//MarshalIndent
		empJSON, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			log.Fatalf(err.Error())
		}

		fmt.Println(string(empJSON))

		// log.("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
	}
}
