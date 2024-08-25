package config

type Config struct {
	Server struct {
		Host           string   `env:"SERVER_HOST, default=0.0.0.0"`
		HttpPort       int      `env:"SERVER_HTTP_PORT, default=8080"`
		HttpsPort      int      `env:"SERVER_HTTPS_PORT, default=8443"`
		TlsCertPath    string   `env:"SERVER_TLS_CERT_PATH"`
		TlsKeyPath     string   `env:"SERVER_TLS_KEY_PATH"`
		TlsCaPaths     []string `env:"SERVER_TLS_CA_PATHS"`
		WriteTimeout   int      `env:"SERVER_WRITE_TIMEOUT, default=15"`
		ReadTimeout    int      `env:"SERVER_READ_TIMEOUT, default=15"`
		IdleTimeout    int      `env:"SERVER_IDLE_TIMEOUT, default=60"`
		TempDir        string   `env:"SERVER_TEMP_DIR, default=."`
		KubeconfigPath string   `env:"KUBECONFIG_PATH, default=~/.kube/config"`
	}

	WebApp struct {
		ContentDir string `env:"APP_CONTENT_DIR, default=./content"`
	}

	// Database struct {
	//     Username string `yaml:"user", env:"DB_USERNAME"`
	//     Password string `yaml:"pass", env:"DB_PASSWORD"`
	// } `yaml:"database"`
}
