package service

import (
	"context"
	"errors"
	"log"
	"time"

	"sprin1/internal/delivery/http/dto"
	"sprin1/internal/model"
)

type reviewService struct {
	reviewRepo ReviewRepository
	userRepo   UserRepository
	// ratingProducer отправляет событие в Kafka для асинхронного пересчёта рейтинга.
	ratingProducer RatingEventProducer
}

// NewReviewService создаёт сервис. Если ratingProducer не nil, пересчёт рейтинга
// будет инициироваться через Kafka; иначе можно добавить fallback-логику.
func NewReviewService(reviewRepo ReviewRepository, userRepo UserRepository, ratingProducer RatingEventProducer) *reviewService {
	return &reviewService{
		reviewRepo:      reviewRepo,
		userRepo:        userRepo,
		ratingProducer:  ratingProducer,
	}
}

func (s *reviewService) CreateReview(ctx context.Context, body dto.CreateReviewRequest) (*model.Review, error) {
	if body.Rating < 1 || body.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}
	if body.AuthorID == body.TargetUserID {
		return nil, errors.New("author_id cannot be equal to target_user_id")
	}
	review := &model.Review{
		DealID:       body.DealID,
		PvzID:        body.PvzID,
		AuthorID:     body.AuthorID,
		TargetUserID: body.TargetUserID,
		Rating:       body.Rating,
		Body:         body.Body,
	}
	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	// После создания отзыва публикуем событие в Kafka для асинхронного пересчёта рейтинга.
	if s.ratingProducer != nil {
		go func(userID int64) {
			// Контекст можно упростить, потому что Producer сам ретраит внутри.
			if err := s.ratingProducer.SendRatingRecalc(context.Background(), userID); err != nil {
				log.Printf("review: failed to send rating recalc event for user %d: %v", userID, err)
			}
		}(body.TargetUserID)
	} else {
		// Fallback: если Kafka не сконфигурирована, пересчитываем рейтинг напрямую, как раньше.
		go func(userID int64) {
			ctx2, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			avg, err := s.reviewRepo.GetAvgRatingByTargetUser(ctx2, userID)
			if err != nil {
				log.Printf("review fallback: get avg for user %d: %v", userID, err)
				return
			}
			if err := s.userRepo.UpdateRatingAvg(ctx2, userID, avg); err != nil {
				log.Printf("review fallback: update rating user %d: %v", userID, err)
			}
		}(body.TargetUserID)
	}

	return review, nil
}

func (s *reviewService) DeleteReview(ctx context.Context, id int64) error {
	return s.reviewRepo.Delete(ctx, id)
}

func (s *reviewService) GetReviews(ctx context.Context, q dto.GetReviewsQuery) ([]*model.Review, int64, error) {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	offset := (q.Page - 1) * q.Limit

	var pvzID *int64
	if q.PvzID != 0 {
		id := q.PvzID
		pvzID = &id
	}

	var targetUserID *int64
	if q.TargetUserID != 0 {
		id := q.TargetUserID
		targetUserID = &id
	}

	filter := ReviewFilter{
		PvzID:        pvzID,
		TargetUserID: targetUserID,
		Offset:       offset,
		Limit:        q.Limit,
	}

	return s.reviewRepo.GetAll(ctx, filter)
}
