package main

import (
	"fmt"
	"sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
	// TODO ініціювати логер
	// TODO ініціювати додаток ( app )
	// TODO запуск gRPC-сервера

}
