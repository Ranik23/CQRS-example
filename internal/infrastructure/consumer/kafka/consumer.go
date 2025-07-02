package kafka

import (
	"context"
	"fmt"
	"log"
	"time"

	"order-service/internal/config"
	"order-service/internal/infrastructure/consumer"

	kafkalib "github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	reader *kafkalib.Reader
}

func (k *kafkaConsumer) Consume(ctx context.Context) (value []byte, key []byte, err error) {
	message, err := k.reader.ReadMessage(ctx)
	if err != nil {
		return nil, nil, err
	}
	return message.Value, message.Key, nil
}

func NewKafkaConsumer(broker string, config *config.Config) (consumer.Consumer, error) {
	var err error
	var conn *kafkalib.Conn

	log.Println("Connecting to Kafka broker:", broker)

	for i := 0; i < 5; i++ {
		conn, err = kafkalib.Dial("tcp", broker)
		if err == nil {
			break
		}
		if i == 4 {
			log.Println("Failed to connect to Kafka after 5 attempts:", err)
			return nil, err
		}
		time.Sleep(5 * time.Second)
	}
	defer conn.Close()

	log.Println("Connected to Kafka broker:", broker)

	controller, err := conn.Controller()
	if err != nil {
		return nil, err
	}

	ctrlConn, err := kafkalib.Dial("tcp", fmt.Sprintf("%s:%d", controller.Host, controller.Port))
	if err != nil {
		return nil, err
	}
	defer ctrlConn.Close()

	topicConfigs := []kafkalib.TopicConfig{
		{
			Topic:             config.Kafka.Topic,
			NumPartitions:     config.Kafka.NumPartitions,
			ReplicationFactor: 1,
		},
	}
	if err := ctrlConn.CreateTopics(topicConfigs...); err != nil {
		return nil, err
	}
	return &kafkaConsumer{
		reader: kafkalib.NewReader(kafkalib.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    config.Kafka.Topic,
			GroupID:  config.Kafka.GroupID, 
			MinBytes: 10e3,          
			MaxBytes: 10e6,       
			MaxWait:  1 * time.Minute,
		}),
	}, nil
}
