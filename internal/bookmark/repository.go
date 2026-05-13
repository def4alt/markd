package bookmark

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, b *Bookmark) error
	Get(ctx context.Context, id string) (*Bookmark, error)
	List(ctx context.Context) ([]Bookmark, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, b *UpdateInput, now time.Time) error
	Search(ctx context.Context, query string) ([]Bookmark, error)
}
