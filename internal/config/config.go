package config

import (
	"github.com/spf13/viper"
	"log"
)

// EnvConfigs Initialize this variable to access the env values
var EnvConfigs *envConfigs

// InitEnvConfigs We will call this in main.go to load the env variables
func InitEnvConfigs() {
	EnvConfigs = loadEnvVariables()
}

// struct to map env values
type envConfigs struct {
	AppServerPort string `mapstructure:"APP_SERVER_PORT"`
	DbDsn         string `mapstructure:"DB_DSN"`
}

// Call to load the variables from env
func loadEnvVariables() (config *envConfigs) {
	// Tell viper the path/location of your env file. If it is root just add "."
	viper.AddConfigPath(".")

	// Tell viper the name of your file
	viper.SetConfigName(".env")

	// Tell viper the type of your file
	viper.SetConfigType("env")

	// Viper reads all the variables from env file and log error if any found
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading env file", err)
	}

	// Viper unmarshal the loaded env variables into the struct
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
	return
}

func init() {
	InitEnvConfigs()
}
