package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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
	var lastCheckedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, q, id).Scan(
		&b.ID,
		&b.URL,
		&b.Title,
		&b.Description,
		&b.Status,
		&b.CreatedAt,
		&b.UpdatedAt,
		&lastCheckedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("bookmark %q not found: %w", id, err)
		}

		return nil, fmt.Errorf("get bookmark %q: %w", id, err)
	}

	if lastCheckedAt.Valid {
		b.LastCheckedAt = &lastCheckedAt.Time
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

		var lastCheckedAt sql.NullTime

		if err := rows.Scan(
			&b.ID,
			&b.URL,
			&b.Title,
			&b.Description,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
			&lastCheckedAt,
		); err != nil {
			return nil, fmt.Errorf("scan bookmark row: %w", err)
		}

		if lastCheckedAt.Valid {
			b.LastCheckedAt = &lastCheckedAt.Time
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

func (r *BookmarkRepository) Update(ctx context.Context, in *bookmark.UpdateInput, now time.Time) error {
	var sets []string
	var args []any

	if in.URL != nil {
		sets = append(sets, "url = ?")
		args = append(args, *in.URL)
	}

	if in.Title != nil {
		sets = append(sets, "title = ?")
		args = append(args, *in.Title)
	}

	if in.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *in.Description)
	}

	if in.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *in.Status)
	}

	if in.LastCheckedAt != nil {
		sets = append(sets, "last_checked_at = ?")
		args = append(args, *in.LastCheckedAt)
	}

	if len(sets) == 0 && in.Tags == nil {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin update bookmark %q tx: %w", in.ID, err)
	}
	defer tx.Rollback()

	if len(sets) > 0 {
		sets = append(sets, "updated_at = ?")
		args = append(args, now, in.ID)

		q := fmt.Sprintf("UPDATE bookmarks SET %s WHERE id = ?", strings.Join(sets, ", "))

		res, err := tx.ExecContext(ctx, q, args...)
		if err != nil {
			return fmt.Errorf("update bookmark %q: %w", in.ID, err)
		}

		if n, _ := res.RowsAffected(); n == 0 {
			return fmt.Errorf("bookmark %q not found", in.ID)
		}
	}

	if in.Tags != nil {
		tx.ExecContext(ctx, `DELETE FROM bookmark_tags WHERE bookmark_id = ?`, in.ID)

		for _, tag := range in.Tags {
			tx.ExecContext(ctx, `INSERT INTO bookmark_tags (bookmark_id, tag) VALUES (?, ?)`, in.ID, tag)
		}
	}

	return tx.Commit()
}

func (r *BookmarkRepository) Create(ctx context.Context, b *bookmark.Bookmark) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin update bookmark %q tx: %w", b.ID, err)
	}
	defer tx.Rollback()

	resp, err := tx.ExecContext(
		ctx,
		`INSERT INTO bookmarks (id, url, title, description, status, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		b.ID, b.URL, b.Title, b.Description, b.Status, b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert bookmark %q: %w", b.ID, err)
	}

	const tagInsertQ = `
		INSERT INTO bookmark_tags (bookmark_id, tag)
		VALUES (?, ?)
	`

	for _, tag := range b.Tags {
		if _, err := tx.ExecContext(ctx, tagInsertQ, b.ID, tag); err != nil {
			return fmt.Errorf("insert bookmark tag %q for %q: %w", tag, b.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create bookmark %q: %w", b.ID, err)
	}

	n, err := resp.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected for bookmark %q: %w", b.ID, err)
	}
	if n == 0 {
		return fmt.Errorf("bookmark %q not created", b.ID)
	}

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

func escapeLikePattern(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `_`, `\_`)

	return s
}

func (r *BookmarkRepository) Search(ctx context.Context, query string) ([]bookmark.Bookmark, error) {
	if query == "" {
		return r.List(ctx)
	}

	const q = `
SELECT id, url, title, description, status, created_at, updated_at, last_checked_at
FROM bookmarks
WHERE
    url LIKE ? ESCAPE '\'
    OR title LIKE ? ESCAPE '\'
    OR description LIKE ? ESCAPE '\'
    OR EXISTS (
        SELECT 1
        FROM bookmark_tags
        WHERE bookmark_id = bookmarks.id
          AND tag LIKE ? ESCAPE '\'
    )
ORDER BY created_at DESC
`

	escaped := escapeLikePattern(query)
	searchQuery := "%" + escaped + "%"

	rows, err := r.db.QueryContext(ctx, q, searchQuery, searchQuery, searchQuery, searchQuery)
	if err != nil {
		return nil, fmt.Errorf("search bookmarks: %w", err)
	}
	defer rows.Close()

	var out []bookmark.Bookmark

	for rows.Next() {
		var b bookmark.Bookmark
		var lastCheckedAt sql.NullTime

		if err := rows.Scan(
			&b.ID,
			&b.URL,
			&b.Title,
			&b.Description,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
			&lastCheckedAt,
		); err != nil {
			return nil, fmt.Errorf("scan bookmark row: %w", err)
		}

		if lastCheckedAt.Valid {
			b.LastCheckedAt = &lastCheckedAt.Time
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
