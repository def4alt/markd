package bookmark

import "time"

type Bookmark struct {
	ID            string
	URL           string
	Title         string
	Description   string
	Status        string
	Tags          []string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastCheckedAt *time.Time
}

type UpdateInput struct {
	ID            string
	URL           *string
	Title         *string
	Description   *string
	Status        *string
	Tags          []string
	LastCheckedAt *time.Time
}
