package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

var (
	// ErrInvalidType defines a invalid type error
	ErrInvalidType = errors.New("invalid model type")
)

// Repository defines interface for model interaction with database
type Repository interface {
	Create(model interface{}) error
	Get(query []Query, model interface{}) error
	GetByID(id int64, model interface{}) error
	GetBetween(start time.Time, end time.Time, model interface{}) error
	DeleteByID(id int64) error
}

// Query defines a valid Repository query
type Query struct {
	Model string
	Field string
	Value interface{}
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	t, _ := q.FormattedQuery()
	fmt.Println(string(t))
	return nil
}
