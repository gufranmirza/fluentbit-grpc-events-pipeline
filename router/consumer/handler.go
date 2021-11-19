package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/api/apiproto"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/encryption"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/kafka"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/pkg/utils"
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
			log.Printf("Failed to unmarshal event: %v", err)
		}

		if c.config.Decrypt {
			c.Decrypt(event)
		}

		if c.config.Print {
			utils.Print(event, c.config.Decrypt)
		}

		err = c.writer.Write(event, c.config.Decrypt)
		if err != nil {
			log.Printf("Failed to write event with error %v", err)
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
