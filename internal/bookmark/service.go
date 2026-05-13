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

type UpdateInput struct {
	ID          string
	URL         string
	Title       string
	Description string
	Tags        []string
}

func (svc *Service) Update(ctx context.Context, in *UpdateInput) error {
	now := svc.clock.Now()

	updated, err := svc.repo.Get(ctx, in.ID)
	if err != nil {
		return err
	}

	updated.UpdatedAt = now

	if in.URL != "" {
		updated.URL = in.URL
	}

	if in.Description != "" {
		updated.Description = in.Description
	}

	if in.Title != "" {
		updated.Title = in.Title
	}

	if in.Tags != nil {
		updated.Tags = in.Tags
	}

	if err := svc.repo.Update(ctx, updated); err != nil {
		return err
	}

	return nil
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

func (svc *Service) Search(ctx context.Context, query string) ([]Bookmark, error) {
	b, err := svc.repo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	return b, nil
}
