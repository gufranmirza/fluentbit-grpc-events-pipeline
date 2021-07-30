# fargo-fb-poc

$ go get google.golang.org/protobuf/cmd/protoc-gen-go
$ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
brew install protobuf

docker build . -t fluentbit-collector -f Dockerfile
docker run -it --rm fluentbit-collector

Server-Side TLS
openssl enc -aes-128-cbc -k secret -P -md sha1
decrypt=YES go run main.go