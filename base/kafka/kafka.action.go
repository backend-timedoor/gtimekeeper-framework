package kafka

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
	kafkaPkg "github.com/segmentio/kafka-go"
)

type Kafka struct {
	SchemaRegistryClient *srclient.SchemaRegistryClient
	Connection           *kafkaPkg.Conn
	Writer               *kafkaPkg.Writer
	Config               *Config
	Brokers              []string
}

type Topic struct {
	Topic       KafkaTopic
	Partition   int
	Replication int
}

type Consumer interface {
	Config() *[]ModuleConfig
}

type Schema struct {
	Subject KafkaTopic
	Type    string
	Schema  string
}

type SchemaRegistry struct {
	Host    string
	Schemas []Schema
}

type Config struct {
	Brokers          string
	Topics           []Topic
	Consumers        []Consumer
	ConsumerGroupID  string
	AutoCommitOffset bool
	SchemaRegistry   SchemaRegistry
}

type ModuleConfig struct {
	Reader kafka.ReaderConfig
	Handle func(context.Context, kafka.Message, *kafka.Reader) error
}

func (k *Kafka) Produce(ctx context.Context, msgs ...kafkaPkg.Message) error {
	errChan := make(chan error, 1)
	var mapNewMessage []kafkaPkg.Message
	for _, msg := range msgs {
		schema, err := k.SchemaRegistryClient.GetLatestSchema(k.getSubject(msg.Topic))
		if err != nil {
			return fmt.Errorf("failed to get latest schema %v", err)
		}

		schemaIDBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(schemaIDBytes, uint32(schema.ID()))

		native, _, err := schema.Codec().NativeFromTextual(msg.Value)
		if err != nil {
			return fmt.Errorf("failed to convert value to native %v", err)
		}

		valueBytes, err := schema.Codec().BinaryFromNative(nil, native)
		if err != nil {
			return fmt.Errorf("failed to convert value to binary %v", err)
		}

		var recordValue []byte
		recordValue = append(recordValue, byte(0))
		recordValue = append(recordValue, schemaIDBytes...)
		recordValue = append(recordValue, valueBytes...)

		mapNewMessage = append(mapNewMessage, kafkaPkg.Message{
			Key:   msg.Key,
			Value: recordValue,
			Topic: msg.Topic,
		})

	}

	go func(errChan chan error) {
		err := k.Writer.WriteMessages(ctx, mapNewMessage...)
		if err != nil {
			errChan <- fmt.Errorf("failed to write messages %v", err)
		}

		close(errChan)
	}(errChan)

	return <-errChan
}

func (k *Kafka) getSubject(topic string) string {
	return topic + "-value"
}
