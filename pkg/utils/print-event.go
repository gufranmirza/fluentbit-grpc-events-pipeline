package utils

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/api/apiproto"
)

func Print(event *apiproto.Event, decrypt bool) {
	e, err := MarshalEvent(event, decrypt)
	if err != nil {
		log.Printf("Failed to print event with error %v", err)
	}
	log.Printf("%s", e)
}

func MarshalEvent(event *apiproto.Event, decrypt bool) (string, error) {
	var err error
	var buf []byte

	if decrypt {
		message := make(map[string]interface{})
		err := json.Unmarshal([]byte(event.Message), &message)
		if err != nil {
			return "", fmt.Errorf("Failed to marshal message %v", err.Error())
		}

		type m struct {
			*apiproto.Event
			Message map[string]interface{} `json:"message"`
		}

		buf, err = json.MarshalIndent(m{Event: event, Message: message}, "", "  ")
		if err != nil {
			return "", fmt.Errorf("Failed to marshal message %v", err.Error())
		}

	} else {
		buf, err = json.MarshalIndent(event, "", "  ")
		if err != nil {
			return "", fmt.Errorf("Failed to marshal message %v", err.Error())
		}
	}

	return string(buf), err
}
