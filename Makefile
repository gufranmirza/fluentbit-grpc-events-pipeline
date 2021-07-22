all: plugin protoc

protoc:
	protoc -I=./api --go_out=./api ./api/api.proto
	protoc -I=./api --go-grpc_out=./api ./api/api.proto

plugin:
	go build -buildmode=c-shared -o ./fluentbit-collector/plugin_grpcout.so ./fluentbit-collector

clean:
	rm -rf *.so *.h *~