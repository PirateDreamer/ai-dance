package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func LoadConfig() (err error) {
	// 获取当前的工作目录
	executable, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return
	}

	// 获取该路径的所在目录
	dir := filepath.Dir(executable)

	viper.SetConfigFile(filepath.Join(dir, "config.yaml"))

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	return
}
