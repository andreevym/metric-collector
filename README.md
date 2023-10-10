# metric-collector

Сервис сбора метрик и алертинга

## Description

Это клиент серверное приложение, где сервер сохраняет и отдает метрики, которые записал клиент. 

# Contents

## [README.md](README.md)
## [cmd](cmd)
- [agent](cmd/agent)
    - [README.md](cmd/agent/README.md)
    - [agent](cmd/agent/agent)
    - [env.go](cmd/agent/env.go)
    - [flag.go](cmd/agent/flag.go)
    - [main.go](cmd/agent/main.go)
- [server](cmd/server)
    - [README.md](cmd/server/README.md)
    - [env.go](cmd/server/env.go)
    - [flag.go](cmd/server/flag.go)
    - [main.go](cmd/server/main.go)
    - [server](cmd/server/server)
## [go.mod](go.mod)
## [go.sum](go.sum)
## [internal](internal)
- [README.md](internal/README.md)
- [counter](internal/counter)
    - [counter.go](internal/counter/counter.go)
    - [counter_test.go](internal/counter/counter_test.go)
- [gauge](internal/gauge)
    - [gauge.go](internal/gauge/gauge.go)
    - [gauge_test.go](internal/gauge/gauge_test.go)
- [handlers](internal/handlers)
    - [get_handler.go](internal/handlers/get_handler.go)
    - [get_handler_test.go](internal/handlers/get_handler_test.go)
    - [handler_test.go](internal/handlers/handler_test.go)
    - [router.go](internal/handlers/router.go)
    - [server.go](internal/handlers/server.go)
    - [update_handler.go](internal/handlers/update_handler.go)
    - [update_handler_test.go](internal/handlers/update_handler_test.go)
- [model](internal/model)
    - [metric.go](internal/model/metric.go)
- [repository](internal/repository)
    - [repository.go](internal/repository/repository.go)
- [server](internal/server)
    - [server.go](internal/server/server.go)
- [storage](internal/storage)
    - [mem](internal/storage/mem)
        - [errors.go](internal/storage/mem/errors.go)
        - [mem.go](internal/storage/mem/mem.go)
        - [mem_test.go](internal/storage/mem/mem_test.go)


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
