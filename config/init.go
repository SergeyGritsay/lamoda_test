package config

import (
	"log"

	"github.com/spf13/viper"
)

const Directory = "/home/sergey/lamoda-test/config"

func getConfigFile() []string {
	return []string{"app-config"}
}

func Init() {
	viper.AddConfigPath(Directory)

	for _, filePath := range getConfigFile() {
		viper.SetConfigName(filePath)
		err := viper.MergeInConfig()
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println(viper.ConfigFileUsed())
}
