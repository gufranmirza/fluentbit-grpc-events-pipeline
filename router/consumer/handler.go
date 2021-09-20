package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/encryption"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/kafka"
	"google.golang.org/protobuf/proto"
)

// Start consumer
func (c *Consumer) Start() {
	consumer := kafka.Consumer{
		Messages: make(chan *sarama.ConsumerMessage, 1000),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.consumer.Consume(ctx, &consumer)
	for message := range consumer.Messages {
		event := &apiproto.Event{}
		err := proto.Unmarshal(message.Value, event)
		if err != nil {
			log.Fatalln("Failed to unmarshal event:", err)
		}

		if c.config.Decrypt {
			c.Decrypt(event)
		}

		if c.config.Print {
			c.Print(event)
		}
	}
}

func (c *Consumer) Decrypt(event *apiproto.Event) {
	key, ok := c.config.AccessTokenDB[event.AccessKey]
	if ok && c.config.Decrypt && key.EncryptionKey != "" {
		msg, err := encryption.Decrypt(string(key.EncryptionKey), string(event.Message))
		if err != nil {
			fmt.Printf("Failed to decrypt message %v/n", err)
		}
		event.Message = fmt.Sprintf("%v", msg)
	}
}

func (c *Consumer) Print(event *apiproto.Event) {
	if c.config.Decrypt {
		message := make(map[string]interface{})
		err := json.Unmarshal([]byte(event.Message), &message)
		if err != nil {
			log.Printf("Failed to marshal message %v", err.Error())
		}

		type m struct {
			*apiproto.Event
			Message map[string]interface{} `json:"message"`
		}

		buf, err := json.MarshalIndent(m{Event: event, Message: message}, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal message %v", err.Error())
		}

		log.Printf("%s", string(buf))
	} else {
		buf, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal message %v", err.Error())
		}

		log.Printf("%s", string(buf))
	}

}
