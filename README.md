# metric-collector

Сервис сбора метрик и алертинга

## Description

Это клиент серверное приложение, где сервер сохраняет и отдает метрики, которые записал клиент. 

## Build Agent

```shell
go build -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=01-01-2024 -X main.buildCommit=05cf15b8f01bf4764d657fd09c7954ea0cdda239" -o agent cmd/agent/main.go  
```

## Build Server

```shell
go build -ldflags "-X main.buildVersion=v1.0.1 -X main.buildDate=01-01-2024 -X main.buildCommit=05cf15b8f01bf4764d657fd09c7954ea0cdda239" -o transport cmd/transport/main.go
```

# Contribution requirements coverage more than 55%

go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out