package producer

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type Producer interface {
	Send(ctx context.Context, key string, payload []byte) error
	Close() error
}

type SaramaProducer struct {
	prod  sarama.SyncProducer
	topic string
}

func NewProducer(brokers []string, topic string, opts ...Option) (*SaramaProducer, error) {
	cfg := PrepareConfig(opts...)

	prod, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create sarama producer: %w", err)
	}

	return &SaramaProducer{
		prod:  prod,
		topic: topic,
	}, nil
}

func (p *SaramaProducer) Send(ctx context.Context, key string, payload []byte) error {
	_ = ctx
	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(key),
		Value:     sarama.ByteEncoder(payload),
		Timestamp: time.Now(),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("pvz"),
				Value: []byte("pvz-sync-prod"),
			},
		},
	}

	_, _, err := p.prod.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("send message failed: %w", err)
	}
	return nil
}

func (p *SaramaProducer) Close() error {
	return p.prod.Close()
}
