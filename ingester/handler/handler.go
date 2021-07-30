package handler

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.ibm.com/Gufran-Baig/fargo-fb-poc/api/apiproto"
)

// Server represents the gRPC server
type Server struct {
	apiproto.UnimplementedEventServiceServer
}

// SayHello generates response to a Ping request
func (s *Server) SendEvent(stream apiproto.EventService_SendEventServer) error {
	// Read Public key for encryption of Events passed over wire
	pubKey, err := ioutil.ReadFile("../cert/encryption_aes.pub")
	if err != nil {
		log.Fatalf("Failed to read key %v \n", err)
	}

	for {
		event, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&apiproto.EResponse{
				Status: apiproto.EventCode_SUCCESS,
			})
		}
		if err != nil {
			return err
		}

		fmt.Println("==============================================")
		if len(os.Getenv("decrypt")) > 0 {
			msg, err := decrypt(string(pubKey), event.Message)
			if err != nil {
				fmt.Printf("Failed to decrypt message %v/n", err)
			}
			event.Message = msg
		}
		fmt.Println(event)

	}
}

func decrypt(key string, ct string) (string, error) {
	ciphertext, err := hex.DecodeString(ct)
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
