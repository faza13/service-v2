package config

type Config struct {
	App      App      `json:"app"`
	Database Database `json:"database"`
	Otel     Otel     `json:"otel"`
	Router   Router   `json:"router"`
	Kafka    Kafka    `json:"kafka"`
	Setting  Setting  `json:"setting"`
	Cache    Cache    `json:"cache"`
	Redis    Redis    `json:"redis"`
	Elastic  Elastic  `json:"elastic"`
}

func NewConfig() Config {
	return Config{
		App: App{
			Name:        "service",
			Environment: "development",
		},
		Database: Database{
			Driver:   "mariadb",
			Host:     "localhost",
			Port:     "3306",
			User:     "root",
			Password: "root",
			Name:     "test",
		},
		Otel: Otel{
			ServiceName: "user",
			HostTempo:   "localhost:4317",
		},
		Router: Router{
			Port:   "9000",
			Prefix: "program",
		},
		Kafka: Kafka{
			Host: "localhost:9092",
		},
		Setting: Setting{
			QueueProgram: "1",
		},
		Cache: Cache{
			Driver: "redis",
		},
		Redis: Redis{
			Name: "0",
			Host: "127.0.0.1:6379",
		},
		Elastic: Elastic{
			Host: "localhost:9200",
		},
	}
}
