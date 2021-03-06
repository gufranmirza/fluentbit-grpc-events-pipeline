# What is does 
This POC is build with the goal to collect events/logs from the host systems such as Kubernetes, docker, VMs etc. 
A buffering layer is added to buffer events from the collector making sure ingestion layer is highly available and collector should be able to dump events without any issues. Events stored in persistance storage e.g. Kafka and routed to destinations as required.

# Design and Arch
Below is the very high level design of the project. It uses below tech stack
- Fluent Bit - As the logs/events collector
- GRPC - As communication layer B/W Fluent bit and  Ingester
- Kafka - Events are stored in Kafka - Persistent storage

![](images/untitled@2x.png)


**Table of Contents**

- [Fluentbit Agent](#fluentbit-agent)
- [Ingester](#ingester)
- [Router](#router)
- [Running The Project](#running-the-project)


# Fluentbit Agent
Fluentbit-Agent is used to collect events from the host system. A Fluentbit plugin is written that communicates with the Ingester API and pushes the events to it.  
- [Plugin Code](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/tree/master/fluentbit-collector/plugin)

When you start the FB-Agent it will require the `ACCESS_KEY` `and ACCESS_TOKEN` 
- `ACCESS_KEY` Configuration for an agent is stored againts this access key. FB-Agent uses access key to obtain the configuration from Ingester API 
- `ACCESS_TOKEN` is the JWT authentication token. FB-Agent will include it with every request it makes to Ingester API. Ingester will validate the token

#### It does the following
- Connect to Ingester
- Exchange the config - requires `ACCESS_TOKEN`
- Stream events to ingester - requires `ACCESS_TOKEN`


#### Config Exchange 
Below is the diagram that shows the config exchange between the FB-Agent and Ingester. Config exchange requires the valid `ACCESS_TOKEN` and `ACCESS_KEY`. It accepst the `ACCESS_KEY` as request parameter.
- You can find the whole [Config Object here](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/blob/master/api/api.proto#L21)
![](images/Config%20Exchange%402x.png)


#### Streaming Events 
It uses the Streaming RPC to stream events from the FB-Agent to the ingester API. It requires the valid `ACCESS_TOKEN` to be passed with every request.
- You can find the whole [Event Object here](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/blob/master/api/api.proto#L8)
![](images/Events@2x.png)


#### Features 
- AES encryption
- Server side TLS 
- Streaming RPC - Events 


#### Open Issues
- Dynamic reloading FB-Agent Config 
- Handling Auth Errors
- Handling Connectivity Error
- Handling Streaming Errors/timeouts 

# Ingester
Ingester is responsible for running the RPC server. FB-Agent sends the events to the Ingester via streaming RPC. Its also has the config for each FB-Agent. It validates the incoming request with passed `ACCESS_TOKEN` and pushes the events to the Kafka topic `fb-kafka`.

- [Ingester Code](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/tree/master/ingester)

#### It does the following 
- Host the config for FB-Agents
- Accept the incoming Events from FB-Agents
- Perform JWT Authentication
- Push Events to Kafka topic `fb-kafka`

NOTE: Temporarily it hosts the config for the FB-Agent. Ideally it should be stored in seperate config service. It loads the config from this [JSON file](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/blob/master/access-tokens-db.json). Where each key is `ACCESS_KEY` that can be used with FB-Agent.

IDEA: Some king of hash can be calculated using customer event data that pushes the events to specific partiton of the kafka topic, and then only one router can be spawn to consume events from that partiton only .

#### Open Issues
- Decouple config management
- Handling Auth Errors
- Handling Kafka Erros

# Router
Router consumes the events from the kafka topic `fb-kafka` and writes it to the specified destinations. Currenlt it is very minimal and it only writes the events to the file (ACCESS_KEY.log) on local diks. 

- [Router Code](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/tree/master/router)

#### It does the following 
- Consume events from `fb-kafka` topic
- Write events to file ACCESS_KEY.log) on local diks. One file is created per ACCESS_KEY

#### Open Issues
- Routing to the different targets
- Ideally you will write the events to COS, S3, Elastic Search etc. Its not implemented yet.

# Running The Project 
### Start Kafka
Go to the project root and run the following command to start the kafka container on docker. Please note that container data is not stored on the disk for now. 
```
docker-compose -f kafka-docker-compose.yaml up -d
```

### Start the ingester
1. Go to the project root and run the following command  to build it
```
cd ingester
go build .
```


2. Run the following command to generate the ACCESS_TOKEN. It will be required when you start the FB-Agent
```
./ingester access-token --expiry=100000000000 
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkzODc1OTk5ODAsImlhdCI6MTYyMTMyMDM0OSwicm9sZSI6IiJ9.8Rpy5M2l-BJ-pD75q8UukLKmSIJvt-O-DytkvtOwbFY
```

Available commands 
```
Usage:
  Ingester access-token [flags]

Flags:
      --expiry int   Expiry duration of token in seconds (default 600)
  -h, --help         help for access-token
```


3. Run the following command to start the Ingester Server 
```
./ingester api --print-events 
```

Available commands
```
Usage:
  Ingester api [flags]

Flags:
  -h, --help           help for api
      --print-events   Print events as received from Fluentbit-Agent
```

### Start the router
1. Open a new terminal, go to the root of the project and run the following command to build it
```
cd router
go build .
```

2. Run the following command to start the router 
```
./router consumer --print-events --decrypt-events
```

Available commands
```
Usage:
  Router consumer [flags]

Flags:
      --decrypt-events   Decrypt events received from kafka - Events written to file also will be decrypted
  -h, --help             help for consumer
      --print-events     Print events as received from Kafka
```

### Start FB-Agent
1. Open a new terminal, go to the root of the project and run the following command to build it
```
docker build . -t fluentbit-collector -f Dockerfile
```

- `ACCESS_TOKEN` - Use the access token you created while starting the ingester
- `ACCESS_KEY` - Use any of the key from this [JSON file](https://github.com/gufranmirza/fluentbit-grpc-events-pipeline/blob/master/access-tokens-db.json)

2. Run following command to run the container image
```
docker run -e ACCESS_KEY=9c60f26f-5b6c-4c80-b5f5-625bf965b6a6 -e ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjkzODc1OTk5ODAsImlhdCI6MTYyMTMyMDM0OSwicm9sZSI6IiJ9.8Rpy5M2l-BJ-pD75q8UukLKmSIJvt-O-DytkvtOwbFY -it --rm fluentbit-collector
```

3. Logs from FluentBit Agent
```
docker run -e ACCESS_KEY=9c60f26f-5b6c-4c80-b5f5-625bf965b6a6 -e ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5c
gufran@Gufran fluentbit-grpc-events-pipeline % docker run -e ACCESS_KEY=9c60f26f-5b6c-4c80-b5f5-625bf965b6a6 -e ACCESS_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjk0MDM1ODYwOTQsImlhdCI6MTYzNzMwNjQ2Miwicm9sZSI6IiJ9.nj4Hr6MJ0F0QbgwVbWoVWtJIYDhHL9tUIfkA93yJcLs -it --rm fluentbit-collector

Fluent Bit v1.7.9
* Copyright (C) 2019-2021 The Fluent Bit Authors
* Copyright (C) 2015-2018 Treasure Data
* Fluent Bit is a CNCF sub-project under the umbrella of Fluentd
* https://fluentbit.io

[2021/11/19 07:52:32] [ info] [engine] started (pid=1)
[2021/11/19 07:52:32] [ info] [storage] version=1.1.1, initializing...
[2021/11/19 07:52:32] [ info] [storage] in-memory
[2021/11/19 07:52:32] [ info] [storage] normal synchronization mode, checksum disabled, max_chunks_up=128
2021/11/19 07:52:32 [out-grpc] plugin parameter = ''
[2021/11/19 07:52:32] [ info] [sp] stream processor started
2021/11/19 07:52:37 4 events sent to ingester
2021/11/19 07:52:37 4 events sent to ingester
2021/11/19 07:52:37 3 events sent to ingester
```
- Observer the terminals to see the activity. Log files ${ACCESS_KEY}.log will be created inside the project root.
- You can spin up more FB-Agent instances with different ACCESS_KEYS

Happy Hacking :-)
