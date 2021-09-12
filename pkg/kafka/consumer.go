package kafka

import (
	"context"
	"log"
	"time"

	"github.com/Shopify/sarama"
)

type ConsumerImpl struct {
	conn          sarama.ConsumerGroup
	Topic         string
	Brokers       []string
	ConsumerGroup string
}

func NewConsumer(topic string, cg string, brokers []string) (*ConsumerImpl, error) {
	c := &ConsumerImpl{
		Topic:         topic,
		ConsumerGroup: cg,
		Brokers:       brokers,
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin // read strategy
	config.Consumer.Offsets.Initial = sarama.OffsetOldest                       // read from begining
	config.Net.DialTimeout = time.Second * 30                                   // Connection dial timeoeout
	config.Net.WriteTimeout = time.Second * 30                                  // Write timeoeout
	config.Net.ReadTimeout = time.Second * 30                                   // Read timeoeout

	/**
	 * Setup a new Sarama consumer group
	 */

	client, err := sarama.NewConsumerGroup(c.Brokers, c.ConsumerGroup, config)
	if err != nil {
		return nil, err
	}
	c.conn = client

	log.Println("Sarama consumer up and running!...")

	return c, nil
}

// call in go routine
func (c *ConsumerImpl) Consume(ctx context.Context, consumer sarama.ConsumerGroupHandler) {
	for {
		// `Consume` should be called inside an infinite loop, when a
		// server-side rebalance happens, the consumer session will need to be
		// recreated to get the new claims
		if err := c.conn.Consume(ctx, []string{c.Topic}, consumer); err != nil {
			log.Panicf("Error from consumer: %v", err)
		}
		// check if context was cancelled, signaling that the consumer should stop
		if ctx.Err() != nil {
			return
		}
	}
}

// close consumer
func (c *ConsumerImpl) Close() error {
	return c.conn.Close()
}
