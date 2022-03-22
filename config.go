package main

type Config struct {
	Server struct {
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
		Port string `yaml:"port", envconfig:"SERVER_PORT"`
	} `yaml:"server"`
	// Database struct {
	//     Username string `yaml:"user", envconfig:"DB_USERNAME"`
	//     Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
	// } `yaml:"database"`
}
