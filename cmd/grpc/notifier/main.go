package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	kafkaBroker = "localhost:9092"
	topic       = "pvz.events-log"
	groupID     = "notifier-group"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaBroker},
		Topic:          topic,
		GroupID:        groupID,
		CommitInterval: time.Second,
	})

	log.Println("[]")

	go func() {
		for {
			msg, err := r.FetchMessage(ctx)
			if err != nil {
				log.Printf("error reading message: %v", err)
				continue
			}
			log.Printf("key: %s, value: %s, partision: %s",
				string(msg.Key),
				string(msg.Value),
				string(rune(msg.Partition)),
			)
			if err := r.CommitMessages(ctx, msg); err != nil {
				log.Printf("error committing message: %v", err)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	log.Println("Graceful shutdown")
	err := r.Close()
	if err != nil {
		return
	}
}
