package config



type KafkaConfig struct {
	Topic   string        `mapstructure:"topics"`
	Brokers string        `mapstructure:"brokers"`
}