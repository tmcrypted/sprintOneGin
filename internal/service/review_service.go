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
	ratingCh   chan int64
}

// NewReviewService создаёт сервис и запускает горутину пересчёта rating_avg по каналу.
func NewReviewService(reviewRepo ReviewRepository, userRepo UserRepository) *reviewService {
	ratingCh := make(chan int64, 100)
	s := &reviewService{reviewRepo: reviewRepo, userRepo: userRepo, ratingCh: ratingCh}
	go s.runRatingUpdater()
	return s
}

func (s *reviewService) runRatingUpdater() {
	for userID := range s.ratingCh {
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		avg, err := s.reviewRepo.GetAvgRatingByTargetUser(ctx, userID)
		if err != nil {
			log.Printf("review service: get avg for user %d: %v", userID, err)
			continue
		}
		if err := s.userRepo.UpdateRatingAvg(ctx, userID, avg); err != nil {
			log.Printf("review service: update rating user %d: %v", userID, err)
		}
	}
}

func (s *reviewService) CreateReview(ctx context.Context, body dto.CreateReviewRequest) (*model.Review, error) {
	if body.Rating < 1 || body.Rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
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
	go func() {
		select {
		case s.ratingCh <- body.TargetUserID:
		default:
			log.Printf("review: rating update queue full, skipping target_user_id=%d", body.TargetUserID)
		}
	}()
	return review, nil
}
