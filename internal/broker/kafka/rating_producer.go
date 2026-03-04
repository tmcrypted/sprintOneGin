package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// RatingRecalcMessage описывает полезную нагрузку события пересчёта рейтинга пользователя.
type RatingRecalcMessage struct {
	TargetUserID int64 `json:"target_user_id"`
}

// RatingProducer отвечает за отправку событий пересчёта рейтинга в Kafka.
type RatingProducer struct {
	producer sarama.SyncProducer
	topic    string
}

// NewRatingProducer создаёт синхронного продьюсера для заданного топика.
func NewRatingProducer(brokers []string, topic string) (*RatingProducer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Return.Successes = true
	cfg.Producer.Retry.Max = 3
	cfg.Producer.Idempotent = true
	cfg.Version = sarama.V3_7_0_0

	p, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		return nil, fmt.Errorf("create kafka producer: %w", err)
	}

	return &RatingProducer{
		producer: p,
		topic:    topic,
	}, nil
}

// Close закрывает продьюсер.
func (p *RatingProducer) Close() error {
	return p.producer.Close()
}

// SendRatingRecalc публикует событие пересчёта рейтинга в Kafka.
func (p *RatingProducer) SendRatingRecalc(_ context.Context, userID int64) error {
	msg := RatingRecalcMessage{
		TargetUserID: userID,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal rating message: %w", err)
	}

	kmsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Value: sarama.ByteEncoder(payload),
		Timestamp: time.Now(),
	}

	if _, _, err := p.producer.SendMessage(kmsg); err != nil {
		return fmt.Errorf("send rating message: %w", err)
	}

	return nil
}

