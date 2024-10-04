package utils

import "github.com/spf13/viper"

type Config struct {
	TmdbAPIKey             string `mapstructure:"TMDB_API_KEY"`
	TmdbAPIReadAccessToken string `mapstructure:"TMDB_API_READ_ACCESS_TOKEN"`
	MongoDbConnectionURI   string `mapstructure:"MONGODB_CONNECTION_URI"`
}

// LoadConfig loads the configuration from the given path
func LoadConfig(path string) (config Config, err error) {
	if path == "" {
		viper.AddConfigPath(".")
	} else {
		viper.AddConfigPath(path)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return

}
