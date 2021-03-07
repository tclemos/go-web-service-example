package services

import (
	"context"

	"github.com/tclemos/go-dockertest-example/internal/core/domain"
	"github.com/tclemos/go-dockertest-example/internal/core/events"
	"github.com/tclemos/go-dockertest-example/internal/core/port"
)

type ThingService struct {
	repo     port.ThingRepository
	notifier port.ThingNotifier
}

func NewThingService(tr port.ThingRepository, tn port.ThingNotifier) *ThingService {
	return &ThingService{
		repo:     tr,
		notifier: tn,
	}
}

func (s *ThingService) Create(ctx context.Context, t domain.Thing) error {
	if err := s.repo.Create(ctx, t); err != nil {
		return err
	}

	tc := events.ThingCreated{Thing: t}
	if err := s.notifier.NotifyThingCreated(tc); err != nil {
		return err
	}

	return nil
}

func (s *ThingService) Get(ctx context.Context, c domain.ThingCode) (*domain.Thing, error) {
	return s.repo.Get(ctx, c)
}
