package config



type KafkaConfig struct {
	Topic   string        `mapstructure:"topics"`
	Brokers string        `mapstructure:"brokers"`
	GroupID string        `mapstructure:"group_id"`
	NumWorkers int           `mapstructure:"num_workers"`
}