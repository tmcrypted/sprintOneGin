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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		avg, err := s.reviewRepo.GetAvgRatingByTargetUser(ctx, userID)
		cancel()
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
	go func() {
		select {
		case s.ratingCh <- body.TargetUserID:
		default:
			log.Printf("review: rating update queue full, skipping target_user_id=%d", body.TargetUserID)
		}
	}()
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
