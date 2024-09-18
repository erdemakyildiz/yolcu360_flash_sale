package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basePath   = filepath.Dir(b)
)

func LoadConfig() {
	env := os.Getenv("profile")

	if env == "test" {
		viper.SetConfigName("test")
	} else {
		viper.SetConfigName("local")
	}

	viper.AddConfigPath(basePath + "/../resource")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
