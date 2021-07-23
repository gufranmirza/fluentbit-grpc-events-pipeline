# fargo-fb-poc

$ go get google.golang.org/protobuf/cmd/protoc-gen-go
$ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
brew install protobuf

docker build . -t fluentbit-collector -f Dockerfile

docker run -it --rm fluentbit-collector

