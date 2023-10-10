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
