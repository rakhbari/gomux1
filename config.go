package main

import (
	"time"
)

type Config struct {
    Server struct {
        Host string `yaml:"host", envconfig:"SERVER_HOST"`
        HttpPort string `yaml:"httpPort", envconfig:"HTTP_PORT"`
		HttpsPort string `yaml:"httpsPort", envconfig:"HTTPS_PORT"`
		CaCertPath string `yaml:"caCertPath", envconfig:"CA_CERT_PATH"`
		CaKeyPath string `yaml:"caKeyPath", envconfig:"CA_KEY_PATH"`
		WriteTimeout time.Duration `yaml:"writeTimeout", envconfig:"WRITE_TIMEOUT"`
		ReadTimeout time.Duration `yaml:"readTimeout", envconfig:"READ_TIMEOUT"`
		IdleTimeout time.Duration `yaml:"idleTimeout", envconfig:"IDLE_TIMEOUT"`
    } `yaml:"server"`
    // Database struct {
    //     Username string `yaml:"user", envconfig:"DB_USERNAME"`
    //     Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
    // } `yaml:"database"`
}
