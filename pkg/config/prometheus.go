package config

type PrometheusConfig struct {
	URL      string `env:"PROMETHEUS_URL,required"`
	Username string `env:"PROMETHEUS_USERNAME"`
	Password string `env:"PROMETHEUS_PASSWORD"`
}
