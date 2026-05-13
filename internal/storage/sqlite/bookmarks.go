package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/def4alt/markd/internal/bookmark"
)

type BookmarkRepository struct {
	db *sql.DB
}

func NewBookmarkRepository(db *sql.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

func (r *BookmarkRepository) Get(ctx context.Context, id string) (*bookmark.Bookmark, error) {
	const q = `
SELECT id, url, title, description, status, created_at, updated_at, last_checked_at
FROM bookmarks
WHERE id = ?
`

	var b bookmark.Bookmark

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&b.ID,
		&b.URL,
		&b.Title,
		&b.Description,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
		&b.LastCheckedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("bookmark %q not found: %w", id, err)
		}

		return nil, fmt.Errorf("get bookmark %q: %w", id, err)
	}

	tags, err := r.listTagsByBookmarkID(ctx, b.ID)
	if err != nil {
		return nil, fmt.Errorf("get bookmark tags %q: %w", id, err)
	}

	b.Tags = tags

	return &b, nil
}

func (r *BookmarkRepository) List(ctx context.Context) ([]bookmark.Bookmark, error) {
	const q = `
SELECT id, url, title, description, status, created_at, updated_at, last_checked_at
FROM bookmarks
ORDER BY created_at DESC
`

	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()

	var out []bookmark.Bookmark

	for rows.Next() {
		var b bookmark.Bookmark

		if err := rows.Scan(
			&b.ID,
			&b.URL,
			&b.Title,
			&b.Description,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
			&b.LastCheckedAt,
		); err != nil {
			return nil, fmt.Errorf("scan bookmark row: %w", err)
		}

		tags, err := r.listTagsByBookmarkID(ctx, b.ID)
		if err != nil {
			return nil, fmt.Errorf("list tags for bookmark %q: %w", b.ID, err)
		}

		b.Tags = tags

		out = append(out, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate bookmarks: %w", err)
	}

	return out, nil
}

func (r *BookmarkRepository) Update(ctx context.Context, b *bookmark.Bookmark) error {
	const q = `
UPDATE bookmarks
SET url = ?, title = ?, description = ?, status = ?, updated_at = ?, last_checked_at = ?
WHERE id = ?
`

	res, err := r.db.ExecContext(
		ctx,
		q,
		b.URL,
		b.Title,
		b.Description,
		b.Status,
		b.UpdatedAt,
		b.ID,
		b.LastCheckedAt,
	)
	if err != nil {
		return fmt.Errorf("update bookmark %q: %w", b.ID, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for bookmark %q: %w", b.ID, err)
	}
	if n == 0 {
		return fmt.Errorf("bookmark %q not found", b.ID)
	}

	return nil
}

func (r *BookmarkRepository) Create(ctx context.Context, b *bookmark.Bookmark) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO bookmarks (id, url, title, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.URL, b.Title, b.Description, b.Status, b.CreatedAt, b.UpdatedAt,
	)

	return err
}

func (r *BookmarkRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin delete bookmark %q tx: %w", id, err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `DELETE FROM bookmark_tags WHERE bookmark_id = ?`, id); err != nil {
		return fmt.Errorf("delete bookmark tags for %q: %w", id, err)
	}

	res, err := tx.ExecContext(ctx, `DELETE FROM bookmarks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete bookmark %q: %w", id, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for delete %q: %w", id, err)
	}
	if n == 0 {
		return fmt.Errorf("bookmark %q not found", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit delete bookmark %q: %w", id, err)
	}

	return nil
}

func (r *BookmarkRepository) listTagsByBookmarkID(ctx context.Context, bookmarkID string) ([]string, error) {
	const q = `
SELECT tag
FROM bookmark_tags
WHERE bookmark_id = ?
ORDER BY tag ASC
`

	rows, err := r.db.QueryContext(ctx, q, bookmarkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string

		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}
