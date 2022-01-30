package configs

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	ServerAddr  string `yaml:"ServerAddr" env:"SERVER_ADDR"`
	MongoPass   string `yaml:"MongoPass" env:"MONGO_PASS"`
	ApiKeyEther string `yaml:"ApiKeyEther" env:"API_ETHER"`
}

var config Config

func GetConfig() (*Config, error) {

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil

	//if err := cleanenv.ReadConfig("config.yml",&config); err != nil {
	//	return nil, err
	//}
	//return &config, nil
}