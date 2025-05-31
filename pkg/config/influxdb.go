package config

type InfluxDB struct {
	Host     string `env:"INFLUXDB_URL,required"`
	Database string `env:"INFLUXDB_DATABASE,required"`
	Token    string `env:"INFLUXDB_TOKEN,required"`
}
