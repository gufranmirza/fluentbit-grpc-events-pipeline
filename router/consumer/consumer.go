package consumer

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/kafka"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/utils"
)

// Config holds the server specific config
type Config struct {
	Decrypt       bool                       // Decrypt events received from kafka
	AccessTokenDB map[string]apiproto.Config // List of access tokens and their config
	Print         bool                       // Print events to console received from kafka
}

// Server represents the gRPC server
type Consumer struct {
	config   *Config
	consumer *kafka.ConsumerImpl
	writer   *utils.Writer
}

// returns consumer impl
func NewConsumer(c *Config) *Consumer {
	consumer := &Consumer{
		config: c,
		writer: utils.New(),
	}

	// read local database of access tokens
	confBytes, err := ioutil.ReadFile("../access-tokens-db.json")
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}
	err = json.Unmarshal(confBytes, &consumer.config.AccessTokenDB)
	if err != nil {
		log.Fatalf("Failed to access tokens db %v \n", err)
	}

	// start kafka consumer
	conn, err := kafka.NewConsumer("plogger-kafka", "plogger-group", []string{"127.0.0.1:9092"})
	if err != nil {
		log.Fatalf("Failed to connect to kafka %v \n", err)
	}
	consumer.consumer = conn

	return consumer
}

func (c *Consumer) CloseConsumer() error {
	c.writer.Close()
	return c.consumer.Close()
}
