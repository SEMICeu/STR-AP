package internal

import (
	"github.com/spf13/viper"
)

func Config() {
	// Initialize Viper
	viper.SetConfigName("env")
	viper.SetConfigFile("./.env")
	viper.AddConfigPath(".")

	// Attempt to read the local configuration
	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		Fatalf("Fatal error config file .env : %s \n", err)
	}
}
