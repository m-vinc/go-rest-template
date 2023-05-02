package models

import (
	"errors"
	"mpj/pkg/ent"

	"entgo.io/ent/dialect/sql"
)

var (
	ErrPaginationLimitNotPositive = errors.New("pagination: limit must be positive")
	ErrPaginationLimitReached     = errors.New("pagination: limit cannot exceed 500")

	ErrPaginationOrderByEmpty    = errors.New("pagination: order_by must be a non-empty string")
	ErrPaginationInvalidOrderDir = errors.New("pagination: order_dir must be equal to 'asc' or 'desc'")

	ErrPaginationOffsetNotPositive = errors.New("pagination: offset must be positive")
)

var ApplicationErrors = map[error]int{
	ErrPaginationLimitNotPositive:  422,
	ErrPaginationLimitReached:      422,
	ErrPaginationOrderByEmpty:      422,
	ErrPaginationInvalidOrderDir:   422,
	ErrPaginationOffsetNotPositive: 422,
}

type PaginationMetadata struct {
	Limit     int `json:"limit"`
	Offset    int `json:"offset"`
	MaxOffset int `json:"max_offset"`
}

type Pagination struct {
	Offset *int64 `form:"offset"`
	Limit  *int64 `form:"limit"`

	OrderBy  *string `form:"order_by"`
	OrderDir *string `form:"order_dir"`
}

func (p *Pagination) Default() *Pagination {
	if p.Offset == nil {
		offset := int64(0)
		p.Offset = &offset
	}

	if p.Limit == nil {
		limit := int64(100)
		p.Limit = &limit
	}

	return p
}

func (p *Pagination) Validate() error {
	if p.Limit != nil {
		if *p.Limit < 1 {
			return ErrPaginationLimitNotPositive
		}
		if *p.Limit > 500 {
			return ErrPaginationLimitReached
		}
	}

	if p.Offset != nil {
		if *p.Offset < 0 {
			return ErrPaginationOffsetNotPositive
		}
	}

	if p.OrderBy != nil {
		if *p.OrderBy == "" {
			return ErrPaginationOrderByEmpty
		}
	}

	if p.OrderDir != nil {
		if *p.OrderDir != "asc" && *p.OrderDir != "desc" {
			return ErrPaginationInvalidOrderDir
		}
	}
	return nil
}

func (pagination *Pagination) Order(defaultBy string, defaultDir func(fields ...string) func(*sql.Selector), allowedBy []string) func(*sql.Selector) {
	by, dir := defaultBy, defaultDir

	if pagination.OrderBy != nil {
		allowed := false
		for _, a := range allowedBy {
			if a == *pagination.OrderBy {
				allowed = true
				break
			}
		}

		if allowed {
			by = *pagination.OrderBy
		}
	}

	if pagination.OrderDir != nil && (*pagination.OrderDir == "asc" || *pagination.OrderDir == "desc") {
		switch *pagination.OrderDir {
		case "asc":
			dir = ent.Asc
		case "desc":
			dir = ent.Desc
		}
	}

	return dir(by)
}
