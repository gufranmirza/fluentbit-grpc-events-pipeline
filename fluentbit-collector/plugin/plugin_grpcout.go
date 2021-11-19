package main

import (
	"context"
	"log"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"github.com/gufranmirza/fluentbit-grpc-events-pipeline/api/apiproto"
	"google.golang.org/grpc"
)
import (
	"C"
	"unsafe"
)

type Plugin struct {
	conn          *grpc.ClientConn
	eventClient   apiproto.EventServiceClient
	config        *apiproto.Config
	encryptionKey string
}

var plugin Plugin

// -------- Plugin Code -------

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "grpc", "GRPC output plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	// init internal plugin object
	plugin = Plugin{config: &apiproto.Config{}}

	// Example to retrieve an optional configuration parameter
	param := output.FLBPluginConfigKey(ctx, "param")
	log.Printf("[out-grpc] plugin parameter = '%s'\n", param)

	// connect to ingester
	err := plugin.connectToIngest()
	if err != nil {
		log.Fatalf("failed to connect err: %v", err)
	}

	// Exchange configuration
	err = plugin.exchangeConfig(plugin.config.AccessKey)
	if err != nil {
		log.Fatalf("failed to exchange config err: %v", err)
	}

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	// Create Fluent Bit decoder
	dec := output.NewDecoder(data, int(length))

	// connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := plugin.eventClient.SendEvent(ctx)
	if err != nil {
		log.Fatalf("failed to connect to stream %v with error %v\n", plugin.conn, err)
	}

	// Iterate Records
	var count int
	for {
		// Record
		ret, ts, event := output.GetRecord(dec)
		if ret != 0 {
			break
		}

		if err := plugin.sendEvent(ts, event, tag, stream); err != nil {
			log.Printf("Error calling sendEvent: %v\n", err)
			continue
		}
		count++
	}

	// receive the response from ingester
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("%v.CloseAndRecv() got error %v, want %v\n", stream, err, nil)
		return output.FLB_RETRY
	}

	if reply.Status == apiproto.EventCode_FAILURE {
		return output.FLB_RETRY
	}

	log.Printf("%v events sent to ingester\n", count)

	// Return options:
	//
	// output.FLB_OK    = data have been processed.
	// output.FLB_ERROR = unrecoverable error, do not try this again.
	// output.FLB_RETRY = retry to flush later.
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	plugin.conn.Close()
	return output.FLB_OK
}

func main() {
}

// ------- Plugin Code ------
