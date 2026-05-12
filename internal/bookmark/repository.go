package bookmark

import "context"

type Repository interface {
	Create(ctx context.Context, b *Bookmark) error
	Get(ctx context.Context, id string) (*Bookmark, error)
	List(ctx context.Context) ([]Bookmark, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, b *Bookmark) error
}
