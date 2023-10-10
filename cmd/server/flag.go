package main

import (
	"flag"
)

// flagRunAddr адрес и порт для запуска сервера, аргумент -a со значением :8080 по умолчанию
var flagRunAddr string

func init() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
}
