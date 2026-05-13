package bookmark

import (
	"context"
	"time"
)

type Clock interface {
	Now() time.Time
}

type IDGenerator interface {
	New() string
}

type Service struct {
	repo  Repository
	clock Clock
	idgen IDGenerator
}

func NewService(repo Repository, clock Clock, idgen IDGenerator) *Service {
	return &Service{
		repo,
		clock,
		idgen,
	}
}

type AddInput struct {
	URL         string
	Title       string
	Description string
	Tags        []string
}

func (svc *Service) Add(ctx context.Context, in AddInput) (*Bookmark, error) {
	now := svc.clock.Now()

	b := &Bookmark{
		ID:          svc.idgen.New(),
		URL:         in.URL,
		Title:       in.Title,
		Description: in.Description,
		Tags:        in.Tags,
		Status:      "unchecked",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := svc.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (svc *Service) Get(ctx context.Context, id string) (*Bookmark, error) {
	b, err := svc.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (svc *Service) Delete(ctx context.Context, id string) error {
	err := svc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (svc *Service) List(ctx context.Context) ([]Bookmark, error) {
	b, err := svc.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	return b, nil
}
