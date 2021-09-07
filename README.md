# fargo-fb-poc

$ go get google.golang.org/protobuf/cmd/protoc-gen-go
$ go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
brew install protobuf

docker build . -t fluentbit-collector -f Dockerfile
docker run -e ACCESS_KEY=9c60f26f-5b6c-4c80-b5f5-625bf965b6a6 -e ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjEwOTQ5NzAsImlhdCI6MTYyMTA4NDk3MCwicm9sZSI6IiJ9.uyGUaazOm5QtpAbfNqY_XwLjcC_WzzV1ll_X8zuZlO0 -it --rm fluentbit-collector

openssl enc -aes-128-cbc -k secret -P -md sha1