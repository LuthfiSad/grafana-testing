package service

import (
	"errors"
	"github.com/user/grafana-analytics-app/internal/models"
	"github.com/user/grafana-analytics-app/internal/repository"
	"time"
)

type ReviewService interface { LeaveReview(productID, customerID uint, rating int, comment string) (*models.Review, error) }
type reviewService struct{ repo repository.ReviewRepository }
func NewReviewService(repo repository.ReviewRepository) ReviewService { return &reviewService{repo} }
func (s *reviewService) LeaveReview(productID, customerID uint, rating int, comment string) (*models.Review, error) {
	if rating < 1 || rating > 5 { return nil, errors.New("invalid rating") }
	rev := &models.Review{ProductID: productID, CustomerID: customerID, Rating: rating, Comment: comment, CreatedAt: time.Now()}
	return rev, s.repo.Create(rev)
}
