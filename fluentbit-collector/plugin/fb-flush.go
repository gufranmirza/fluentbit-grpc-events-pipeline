package main

import (
	"encoding/json"
	"os"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/pkg/encryption"
	"google.golang.org/protobuf/types/known/timestamppb"
)
import (
	"C"
	"fmt"
)

func (plugin *Plugin) sendEvent(ts interface{}, event map[interface{}]interface{}, tag *C.char, stream apiproto.EventService_SendEventClient) error {
	// Timestamp
	var timestamp time.Time
	switch tts := ts.(type) {
	case output.FLBTime:
		timestamp = tts.Time
	case uint64:
		// From our observation, when ts is of type uint64 it appears to
		// be the amount of seconds since unix epoch.
		timestamp = time.Unix(int64(tts), 0)
	default:
		timestamp = time.Now()
	}

	var (
		encryptedEvent string
		err            error
	)

	ev := make(map[string]interface{})
	for key, value := range event {
		strKey := fmt.Sprintf("%v", key)
		ev[strKey] = value
	}

	buffer, err := json.Marshal(ev)
	if err != nil {
		return err
	}

	if plugin.config.EncryptionKey != "" {
		encryptedEvent, err = encryption.Encrypt(plugin.config.EncryptionKey, string(buffer))
		if err != nil {
			return fmt.Errorf("failed to encrpt message %v", err)
		}
	} else {
		encryptedEvent = string(buffer)
	}

	hostname, _ := os.Hostname()
	e := &apiproto.Event{
		Tag:     C.GoString(tag),
		AgentId: hostname,
		Message: encryptedEvent,
		Timestamp: &timestamppb.Timestamp{
			Seconds: timestamp.UnixNano(),
			Nanos:   int32(timestamp.UnixNano()),
		},
		UserID:    plugin.config.UserID,
		AccessKey: plugin.config.AccessKey,
	}

	if err := stream.Send(e); err != nil {
		return fmt.Errorf("failed to send event with error %v", err)
	}

	return nil
}
