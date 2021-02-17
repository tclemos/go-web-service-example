package services

import (
	"context"

	"github.com/tclemos/go-dockertest-example/internal/core/domain"
	"github.com/tclemos/go-dockertest-example/internal/core/port"
)

type ThingService struct {
	repo port.ThingRepository
}

func NewThingService(tr port.ThingRepository) *ThingService {
	return &ThingService{
		repo: tr,
	}
}

func (s *ThingService) Create(ctx context.Context, t domain.Thing) error {
	return s.repo.Create(ctx, t)
}

func (s *ThingService) Get(ctx context.Context, c domain.ThingCode) (*domain.Thing, error) {
	return s.repo.Get(ctx, c)
}
