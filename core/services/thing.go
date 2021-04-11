package services

import (
	"context"

	"github.com/tclemos/go-web-service-example/core/domain"
	"github.com/tclemos/go-web-service-example/core/ports"
)

type ThingService struct {
	repo     ports.ThingRepository
	notifier ports.ThingNotifier
}

func NewThingService(tr ports.ThingRepository, tn ports.ThingNotifier) *ThingService {
	return &ThingService{
		repo:     tr,
		notifier: tn,
	}
}

func (s *ThingService) Create(ctx context.Context, t domain.Thing) error {
	if err := s.repo.Create(ctx, t); err != nil {
		return err
	}

	tc := domain.ThingCreated{Thing: t}
	if err := s.notifier.NotifyThingCreated(tc); err != nil {
		return err
	}

	return nil
}

func (s *ThingService) Get(ctx context.Context, id domain.ThingID) (*domain.Thing, error) {
	return s.repo.Get(ctx, id)
}

func (s *ThingService) GetByCode(ctx context.Context, c domain.ThingCode) (*domain.Thing, error) {
	return s.repo.GetByCode(ctx, c)
}
