package adapters

import (
	"github.com/IBM/sarama"
)

type KafkaConsumer struct {
	consumer sarama.Consumer
}

func NewKafkaConsumer(brokers []string) (*KafkaConsumer, error) {
	consumer, err := sarama.NewConsumer(brokers, nil)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{consumer: consumer}, nil
}

func (kc *KafkaConsumer) Consume(topic string, partition int32, messageHandler func(*sarama.ConsumerMessage)) error {
	partitionConsumer, err := kc.consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if err != nil {
		return err
	}
	defer partitionConsumer.Close()

	// Process messages
	for message := range partitionConsumer.Messages() {
		messageHandler(message)
	}

	return nil
}

func (kc *KafkaConsumer) Close() error {
	return kc.consumer.Close()
}
