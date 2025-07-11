package kafka_test

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/require"
)

const (
	saramaTopic  = "pvz.events-log-test"
	saramaBroker = "localhost:9092"
)

func Test_SendAndReceive_WithConsumerGroup(t *testing.T) {
	clusterAdmin, err := sarama.NewClusterAdmin([]string{saramaBroker}, sarama.NewConfig())
	require.NoError(t, err)

	defer func(admin sarama.ClusterAdmin) {
		err := admin.Close()
		if err != nil {
			_ = fmt.Errorf("error closing cluster clusterAdmin: %v", err)
		}
	}(clusterAdmin)

	_ = clusterAdmin.DeleteTopic(saramaTopic)
	err = clusterAdmin.CreateTopic(saramaTopic, &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
	}, false)
	require.NoError(t, err)

	payload := map[string]string{
		"event_id":   "123",
		"event_type": "order_accepted",
	}
	data, err := json.Marshal(payload)
	require.NoError(t, err)

	prodCfg := sarama.NewConfig()
	prodCfg.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{saramaBroker}, prodCfg)
	require.NoError(t, err)

	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: saramaTopic,
		Key:   sarama.StringEncoder("integration"),
		Value: sarama.ByteEncoder(data),
	})
	require.NoError(t, err)

	err = producer.Close()
	if err != nil {
		_ = fmt.Errorf("error closing producer: %v", err)
	}
	consumerCfg := sarama.NewConfig()
	consumerCfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	consumerCfg.Consumer.Offsets.Initial = sarama.OffsetOldest

	groupID := "pvz-events-log-test-group"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	consumerGroup, err := sarama.NewConsumerGroup([]string{saramaBroker}, groupID, consumerCfg)
	require.NoError(t, err)
	defer func(consumerGroup sarama.ConsumerGroup) {
		err := consumerGroup.Close()
		if err != nil {
			_ = fmt.Errorf("error closing consumer group: %v", err)
		}
	}(consumerGroup)

	handler := &testConsumerGroupHandler{
		expectedKey:   "integration",
		expectedValue: payload,
		test:          t,
		done:          make(chan struct{}),
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			if err := consumerGroup.Consume(ctx, []string{saramaTopic}, handler); err != nil {
				t.Logf("error from consumer: %v", err)
				return
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	select {
	case <-handler.done:

	case <-ctx.Done():
		t.Fatal("timeout waiting for message")
	}

	wg.Wait()
}

type testConsumerGroupHandler struct {
	expectedKey   string
	expectedValue map[string]string
	test          *testing.T
	done          chan struct{}
}

func (h *testConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *testConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *testConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		require.Equal(h.test, h.expectedKey, string(msg.Key))
		var received map[string]string
		err := json.Unmarshal(msg.Value, &received)
		require.NoError(h.test, err)
		require.Equal(h.test, h.expectedValue, received)

		sess.MarkMessage(msg, "")
		close(h.done)
		break
	}
	return nil
}
