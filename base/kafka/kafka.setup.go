package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/backend-timedoor/gtimekeeper/app"
	"github.com/backend-timedoor/gtimekeeper/base/contracts"
	"github.com/segmentio/kafka-go"
)


func BootKafka(topics []string, consumers []contracts.KafkaConsumer) contracts.Kafka {
	fmt.Println("initiate kafka...")
	// kafka setup
	stringBroker := app.Config.GetString("kafka.brokers")
	brokers := strings.Split(stringBroker, ",")

	// kafka create topic
	createTopic(topics, brokers)

	// kafka consumer register
	initConsumer(consumers, brokers)

	fmt.Println("kafka ready to fire...")
	return &Kafka{
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  brokers,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func createTopic(topics []string, brokers []string) {
	for _, topic := range topics {
		for _, broker := range brokers {
			con, err := kafka.DialLeader(context.Background(), "tcp", broker, topic, 0)
			if err != nil {
				app.Log.Fatalf("cannot register topic kafka %s broker %s: %v", topic, broker, err)
			}

			con.Close()
		}

	}

	fmt.Println("finish register topics kafka")
}

func initConsumer(consumers []contracts.KafkaConsumer, brokers []string) {
	for _, consumer := range consumers {
		go func() {
			reader := kafka.NewReader(kafka.ReaderConfig{
				Brokers:   brokers,
				GroupID:   consumer.Group(),
				Topic:     consumer.Topic(),
			})
		
			for {
				m, err := reader.ReadMessage(context.Background())
				if err != nil {
					break
				}
				consumer.Handle(m)
				// fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
			}
		}()
	}

	fmt.Println("finish register consumers kafka")
}