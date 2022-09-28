package config

type Config struct {
    Server struct {
        Host string `env:"SERVER_HOST, default=0.0.0.0"`
        HttpPort int `env:"SERVER_HTTP_PORT, default=8080"`
        HttpsPort int `env:"SERVER_HTTPS_PORT, default=8443"`
        TlsCertPath string `env:"SERVER_TLS_CERT_PATH, require=true"`
        TlsKeyPath string `env:"SERVER_TLS_KEY_PATH, require=true"`
        TlsCaPaths []string `env:"SERVER_TLS_CA_PATHS"`
        WriteTimeout int `env:"SERVER_WRITE_TIMEOUT, default=15"`
        ReadTimeout int `env:"SERVER_READ_TIMEOUT, default=15"`
        IdleTimeout int `env:"SERVER_IDLE_TIMEOUT, default=60"`
    }
    
    // Database struct {
    //     Username string `yaml:"user", env:"DB_USERNAME"`
    //     Password string `yaml:"pass", env:"DB_PASSWORD"`
    // } `yaml:"database"`
}