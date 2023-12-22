# metric-collector

Сервис сбора метрик и алертинга

## Description

Это клиент серверное приложение, где сервер сохраняет и отдает метрики, которые записал клиент. 

## Build Agent

```shell
go build -o agent cmd/agent/main.go  
```

## Build Server

```shell
go build -o server cmd/server/main.go
```
