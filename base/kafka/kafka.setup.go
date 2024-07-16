package kafka

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/backend-timedoor/gtimekeeper-framework/container"
	"github.com/backend-timedoor/gtimekeeper-framework/utils/helper"
	"github.com/riferrei/srclient"
	"github.com/segmentio/kafka-go"
)

const ContainerName string = "kafka"

type KafkaTopic string

func New(config *Config) *Kafka {
	fmt.Println("initiate kafka...")
	// kafka setup
	brokers := strings.Split(config.Brokers, ",")

	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		log.Fatalf("failed to dial leader: %v", err)
	}
	defer conn.Close()

	k := &Kafka{
		SchemaRegistryClient: srclient.CreateSchemaRegistryClient(config.SchemaRegistry.Host),
		Connection:           conn,
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
		}),
		Config:  config,
		Brokers: brokers,
	}

	k.initTopic(config.Topics)
	k.initSchemaRegistry(config.SchemaRegistry)
	k.initConsumer(config.Consumers)

	container.Set(ContainerName, k)

	return k
}

func (k *Kafka) initTopic(topics *[]Topic) {
	var topicConfigs []kafka.TopicConfig

	for _, t := range *topics {
		partition := 1
		replication := 1

		if t.Partition != 0 || t.Partition < 0 {
			partition = t.Partition
		}

		if t.Replication != 0 || t.Replication < 0 {
			replication = t.Replication
		}

		topicConfigs = append(topicConfigs, kafka.TopicConfig{
			Topic:             string(t.Topic),
			NumPartitions:     partition,
			ReplicationFactor: replication,
		})
	}

	k.createTopic(topicConfigs)

	fmt.Println("finish register topics kafka")
}

func (k *Kafka) initSchemaRegistry(schemaRegistry SchemaRegistry) {
	for _, s := range *schemaRegistry.Schemas {
		schemaBytes, err := os.ReadFile(s.Schema)
		if err != nil {
			log.Fatalf("failed to read schema file: %v", err)
		}

		subject := string(s.Subject + "-value")
		_, err = k.SchemaRegistryClient.CreateSchema(
			string(subject),
			string(schemaBytes),
			srclient.Avro,
		)
		if err != nil {
			log.Fatalf("failed to create schema: %v", err)
		}

		_, err = k.SchemaRegistryClient.ChangeSubjectCompatibilityLevel(subject, srclient.None)
		if err != nil {
			log.Fatalf("failed to change subject compatibility level: %v", err)
		}
	}

	fmt.Println("finish register schema kafka")
}

func (k *Kafka) initConsumer(consumers *[]Consumer) {
	for _, consumer := range *consumers {
		configs := consumer.Config()

		for _, config := range *configs {
			var readerConfig kafka.ReaderConfig
			helper.Clone(&readerConfig, &config.Reader)

			if k.Brokers != nil {
				readerConfig.Brokers = k.Brokers
			}

			if k.Config.ConsumerGroupID != "" {
				readerConfig.GroupID = k.Config.ConsumerGroupID
			}

			go k.readMessage(config, readerConfig)
		}
	}

	fmt.Println("finish register consumers kafka")
}

func (k *Kafka) readMessage(consumerConfig ModuleConfig, readerConfig kafka.ReaderConfig) {
	ctx := context.Background()
	reader := kafka.NewReader(readerConfig)

	defer reader.Close()

	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			fmt.Printf("error fetch message kafka: %v\n", err)
			break
		}

		schema, err := k.SchemaRegistryClient.GetLatestSchema(k.getSubject(consumerConfig.Reader.Topic))
		if err != nil {
			log.Fatalf("failed to get latest schema %v", err)
		}
		native, _, _ := schema.Codec().NativeFromBinary(m.Value[5:])
		value, _ := schema.Codec().TextualFromNative(nil, native)

		m.Value = value

		err = consumerConfig.Handle(ctx, m, reader)
		if err != nil {
			fmt.Printf("error handling: %v\n", err)
		} else if k.Config.AutoCommitOffset {
			if err := reader.CommitMessages(context.Background(), m); err != nil {
				fmt.Printf("error commit message: %v\n", err)
			}
		}
	}
}

func (k *Kafka) createTopic(topicConfigs []kafka.TopicConfig) {
	controller, err := k.Connection.Controller()
	if err != nil {
		log.Fatalf(err.Error())
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer controllerConn.Close()

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		log.Fatalf(err.Error())
	}
}
