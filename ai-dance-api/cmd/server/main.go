package main

import (
	"ai-dance-api/config"
	"ai-dance-api/internal"

	"go.uber.org/fx"
)

func main() {
	app := fx.New()

	conf := config.InitConfigs()
	internal.InitHttp(conf)
	app.Run()
}
