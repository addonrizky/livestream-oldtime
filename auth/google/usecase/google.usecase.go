package usecase

import (
	"context"
	"time"

	"github.com/asumsi/livestream/auth/google"
)

type googleUsecase struct {
	googleRepo     google.GoogleRepository
	contextTimeout time.Duration
}

func NewGoogleUsecase(googleRepo google.GoogleRepository, timeout time.Duration) google.GoogleUsecase {
	return &googleUsecase{
		googleRepo:     googleRepo,
		contextTimeout: timeout,
	}
}

func (a *googleUsecase) CheckIsUserAlreadyRegistered(ctx context.Context, email string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	count, err := a.googleRepo.CountGoogleUserByEmail(ctx, email)
	if count < 1 || err != nil {
		return false, err
	}

	return true, err
}
