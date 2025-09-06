package main

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadConfig() (err error) {
	viper.SetConfigFile("D:/workspace/mcp-project/mcp-mysql/server/config.yaml")

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	return
}
