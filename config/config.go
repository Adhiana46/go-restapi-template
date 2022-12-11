package config

type Config struct {
	Host string `env:"HOST" env-default:""`
	Port string `env:"PORT" env-default:"8000"`

	// Database
	DbHost    string `env:"DB_HOST" env-default:"localhost"`
	DbPort    string `env:"DB_PORT" env-default:"5432"`
	DbUser    string `env:"DB_USER" env-default:"user"`
	DbPass    string `env:"DB_PASS" env-default:"secret"`
	DbName    string `env:"DB_NAME" env-default:"todoapp"`
	DbSSL     string `env:"DB_SSL" env-default:"disable"`
	DbDialect string `env:"DB_DIALECT" env-default:"pgx"`

	// RabbitMQ
	AmqpHost string `env:"AMQP_HOST" env-default:"localhost"`
	AmqpPort string `env:"AMQP_PORT" env-default:"5672"`
	AmqpUser string `env:"AMQP_USER" env-default:"guest"`
	AmqpPass string `env:"AMQP_PASS" env-default:"guest"`
}
