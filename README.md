# metric-collector

Сервис сбора метрик и алертинга

## Description

Это клиент серверное приложение, где сервер сохраняет и отдает метрики, которые записал клиент. 

## Build Agent

```shell
cd $GOPATH/src/github.com/andreevym/metric-collector/cmd/agent
go build -o agent *.go
```

## Build Server

```shell
cd $GOPATH/src/github.com/andreevym/metric-collector/cmd/server
go build -o server *.go
```

SERVER_PORT=$(random unused-port)
ADDRESS="localhost:${SERVER_PORT}"
TEMP_FILE=$(random tempfile)
metricstest -test.v -test.run=^TestIteration4$ \
-agent-binary-path=cmd/agent/agent \
-binary-path=cmd/server/server \
-server-port=$SERVER_PORT \
-source-path=.