package kafka

import (
	"context"
	"fmt"

	"github.com/backend-timedoor/gtimekeeper-framework/app"
	kafkaPkg "github.com/segmentio/kafka-go"
)


type Kafka struct {
	Writer *kafkaPkg.Writer
}

func (k *Kafka) Produce(msgs ...kafkaPkg.Message) {
	fmt.Println("messages register %v", len(msgs))
	err := k.Writer.WriteMessages(context.Background(), msgs...)

	if err != nil {
		app.Log.Fatalf("failed to write messages %v", err)
	}
}