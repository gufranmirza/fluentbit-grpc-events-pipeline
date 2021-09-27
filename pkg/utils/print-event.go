package utils

import (
	"encoding/json"
	"log"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
)

func Print(event *apiproto.Event, decrypt bool) {
	if decrypt {
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
