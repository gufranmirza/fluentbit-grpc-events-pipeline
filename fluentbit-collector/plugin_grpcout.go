package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fluent/fluent-bit-go/output"
	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/types/known/timestamppb"
)
import (
	"C"
	"fmt"
	"unsafe"
)

var clientConn *grpc.ClientConn
var c apiproto.EventServiceClient
var encryptionKey string

//export FLBPluginRegister
func FLBPluginRegister(ctx unsafe.Pointer) int {
	return output.FLBPluginRegister(ctx, "grpc", "GRPC output plugin")
}

//export FLBPluginInit
func FLBPluginInit(ctx unsafe.Pointer) int {
	// Example to retrieve an optional configuration parameter
	param := output.FLBPluginConfigKey(ctx, "param")
	fmt.Printf("[out-grpc] plugin parameter = '%s'\n", param)

	var err error

	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile("/fluent-bit/bin/ca-cert.pem", "x.test.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Dial
	clientConn, err = grpc.Dial("host.docker.internal:7777", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}

	// setup streaming
	c = apiproto.NewEventServiceClient(clientConn)

	// Read Public key for encryption of Events passed over wire
	pubKey, err := ioutil.ReadFile("/fluent-bit/bin/encryption_aes.pub")
	if err != nil {
		log.Fatalf("Failed to read key %v \n", err)
	}
	encryptionKey = string(pubKey)

	return output.FLB_OK
}

//export FLBPluginFlush
func FLBPluginFlush(data unsafe.Pointer, length C.int, tag *C.char) int {
	var (
		ret    int
		ts     interface{}
		record map[interface{}]interface{}
	)

	// Create Fluent Bit decoder
	dec := output.NewDecoder(data, int(length))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.SendEvent(ctx)
	if err != nil {
		log.Fatalf("%v.SendEvent(_) = _, %v", c, err)
	}

	// Iterate Records
	for {
		// Record
		ret, ts, record = output.GetRecord(dec)
		if ret != 0 {
			break
		}

		rec := make(map[string]string)
		for k, v := range record {
			strKey := fmt.Sprintf("%s", k)
			strValue := fmt.Sprintf("%s", v)
			rec[strKey] = strValue
		}

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

		encryptedEvent, err := encrypt(encryptionKey, fmt.Sprintf("%v", rec))
		if err != nil {
			log.Printf("Failed to encrpt message %v", err)
			continue
		}

		hostname, _ := os.Hostname()
		event := &apiproto.Event{
			Tag:     C.GoString(tag),
			AgentId: hostname,
			Message: encryptedEvent,
			Timestamp: &timestamppb.Timestamp{
				Seconds: timestamp.UnixNano(),
				Nanos:   int32(timestamp.UnixNano()),
			},
		}
		log.Printf("%v\n", record)

		if err := stream.Send(event); err != nil {
			log.Fatalf("Error calling RecordEvents: %s", err)
		}
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)

	if reply.Status == apiproto.EventCode_FAILURE {
		return output.FLB_RETRY
	}

	// Return options:
	//
	// output.FLB_OK    = data have been processed.
	// output.FLB_ERROR = unrecoverable error, do not try this again.
	// output.FLB_RETRY = retry to flush later.
	return output.FLB_OK
}

//export FLBPluginExit
func FLBPluginExit() int {
	clientConn.Close()
	return output.FLB_OK
}

func main() {
}

// ----- Non Plugin Code -------
func encrypt(key string, message string) (string, error) {
	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher([]byte(key))
	// if there are any errors, handle them
	if err != nil {
		return "", err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		fmt.Println(err)
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		fmt.Println(err)
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	ciphertext := gcm.Seal(nonce, nonce, []byte(message), nil)
	return hex.EncodeToString(ciphertext), nil
}
