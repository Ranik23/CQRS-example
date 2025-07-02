package kafka

import (
	"context"
	"order-service/internal/infrastructure/producer"
	kafkalib "github.com/segmentio/kafka-go"
)

type kafkaProducer struct {
	writer *kafkalib.Writer
}

// Produce implements producer.Producer.
func (k *kafkaProducer) Produce(ctx context.Context, key []byte, message []byte) error {
	if err := k.writer.WriteMessages(ctx, kafkalib.Message{
		Key: key,
		Value: message,
	}); err != nil {
		return err
	}
	return nil
}

func NewKafkaProducer(brokers string, topic string) producer.Producer {
	return &kafkaProducer{
		writer: kafkalib.NewWriter(kafkalib.WriterConfig{
			Brokers: []string{brokers},
			Topic: topic,
		}),
	}
}
