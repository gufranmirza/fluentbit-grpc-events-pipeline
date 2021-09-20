package kafka

import (
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

type Producer struct {
	conn    sarama.AsyncProducer
	Topic   string
	Brokers []string
}

func NewProducer(topic string, brokers []string) (*Producer, error) {
	p := &Producer{
		Topic:   topic,
		Brokers: brokers,
	}

	config := sarama.NewConfig()

	config.Producer.Return.Errors = true                   // Return producer error messages
	config.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.MaxMessageBytes = 8 * 1024 * 1024

	config.Producer.Retry.Backoff = time.Second * 30 // Next retry wait time
	config.Producer.Retry.Max = 5                    // Max retries to produce message
	config.Net.DialTimeout = time.Second * 30        // Connection dial timeoeout
	config.Net.WriteTimeout = time.Second * 30       // Write timeoeout
	config.Net.ReadTimeout = time.Second * 30        // Read timeoeout

	hostname, _ := os.Hostname()
	config.ClientID = hostname // Set hostname as clientID

	conn, err := sarama.NewAsyncProducer(p.Brokers, config)
	if err != nil {
		return nil, err
	}
	p.conn = conn

	go func() {
		for err := range p.conn.Errors() {
			log.Printf("[ERROR] - Producer %v \n", err)
		}
	}()

	return p, nil

}

func (p *Producer) Produce(message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: p.Topic,
		Value: sarama.ByteEncoder(message),
	}

	p.conn.Input() <- msg
}

func (p *Producer) Close() error {
	return p.conn.Close()
}
