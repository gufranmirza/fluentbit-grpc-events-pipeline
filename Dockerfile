FROM golang:1.16 as gobuilder

WORKDIR /root
COPY . /root/

ENV GOOS=linux\
    GOARCH=amd64\
    GO111MODULE=on

RUN go mod download 
RUN go build -buildmode=c-shared -o plugin_grpcout.so ./fluentbit-collector/plugin

FROM fluent/fluent-bit:1.7

COPY --from=gobuilder /root/plugin_grpcout.so /fluent-bit/bin/
COPY --from=gobuilder /root/cert/ca-cert.pem /fluent-bit/bin/
COPY --from=gobuilder /root/fluentbit-collector/fluent-bit.conf /fluent-bit/etc/
COPY --from=gobuilder /root/fluentbit-collector/plugins.conf /fluent-bit/etc/

EXPOSE 2020

CMD ["/fluent-bit/bin/fluent-bit", "--config", "/fluent-bit/etc/fluent-bit.conf"]
