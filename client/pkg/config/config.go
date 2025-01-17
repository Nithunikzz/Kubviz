package config

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
	NatsToken   string `envconfig:"NATS_TOKEN"`
	DbPort      int    `envconfig:"DB_PORT"`
	DBAddress   string `envconfig:"DB_ADDRESS"`
}
