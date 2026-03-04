package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"sprin1/internal/service"

	"github.com/IBM/sarama"
)

// RatingRecalcConsumer обрабатывает события пересчёта рейтинга из Kafka
// и вызывает репозитории домена для обновления среднего рейтинга пользователя.
type RatingRecalcConsumer struct {
	reviewRepo service.ReviewRepository
	userRepo   service.UserRepository
}

// NewRatingRecalcConsumer создаёт consumer group handler для Sarama.
func NewRatingRecalcConsumer(reviewRepo service.ReviewRepository, userRepo service.UserRepository) *RatingRecalcConsumer {
	return &RatingRecalcConsumer{
		reviewRepo: reviewRepo,
		userRepo:   userRepo,
	}
}

// Setup вызывается при запуске consumer group.
func (h *RatingRecalcConsumer) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup вызывается при остановке consumer group.
func (h *RatingRecalcConsumer) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim читает сообщения из партиции и пересчитывает рейтинг.
func (h *RatingRecalcConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var payload RatingRecalcMessage
		if err := json.Unmarshal(msg.Value, &payload); err != nil {
			log.Printf("kafka rating consumer: invalid payload: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		avg, err := h.reviewRepo.GetAvgRatingByTargetUser(ctx, payload.TargetUserID)
		cancel()
		if err != nil {
			log.Printf("kafka rating consumer: get avg for user %d: %v", payload.TargetUserID, err)
			session.MarkMessage(msg, "")
			continue
		}

		ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
		if err := h.userRepo.UpdateRatingAvg(ctx2, payload.TargetUserID, avg); err != nil {
			log.Printf("kafka rating consumer: update rating for user %d: %v", payload.TargetUserID, err)
		}
		cancel2()

		session.MarkMessage(msg, "")
	}

	return nil
}

// RunRatingConsumer запускает consumer group в отдельной горутине и
// корректно останавливается по сигналу завершения.
func RunRatingConsumer(brokers []string, groupID, topic string, reviewRepo service.ReviewRepository, userRepo service.UserRepository) error {
	cfg := sarama.NewConfig()
	cfg.Version = sarama.V3_7_0_0
	cfg.Consumer.Offsets.Initial = sarama.OffsetNewest

	group, err := sarama.NewConsumerGroup(brokers, groupID, cfg)
	if err != nil {
		return err
	}

	handler := NewRatingRecalcConsumer(reviewRepo, userRepo)

	ctx, cancel := context.WithCancel(context.Background())

	// Завершение по SIGINT/SIGTERM.
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	go func() {
		defer group.Close()
		for {
			if err := group.Consume(ctx, []string{topic}, handler); err != nil {
				log.Printf("kafka rating consumer: consume error: %v", err)
				time.Sleep(2 * time.Second)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return nil
}

