package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/booking/database"
	"github.com/booking/model"
)

var (
	// ErrRoomDNE defined a Room does not exist error
	ErrRoomDNE error = errors.New("room does not exist")
	// ErrRoomExistsError defined a Room already exists
	ErrRoomExistsError error = errors.New("room already exist")
)

type roomRepository struct {
	db database.Database
}

// NewRoomRepository returns a room implementation of Repository
func NewRoomRepository(db database.Database, log bool) (Repository, error) {
	if err := db.CreateSchema([]interface{}{
		(*model.Room)(nil),
	}); err != nil {
		return nil, err
	}

	if log {
		db.Conn().AddQueryHook(dbLogger{})
	}

	return &roomRepository{
		db: db,
	}, nil
}

func (r *roomRepository) Create(m interface{}) error {
	room, ok := m.(*model.Room)
	if !ok {
		return ErrInvalidType
	}
	_, err := r.db.Conn().Model(room).Insert()
	return roomError(err)
}

func (r *roomRepository) Get(q []Query, m interface{}) error {
	rooms, ok := m.(*[]model.Room)
	if !ok {
		return ErrInvalidType
	}

	query := r.db.Conn().Model(rooms)

	for _, v := range q {
		query = query.Where(fmt.Sprintf("%v.%v = ?", v.Model, v.Field), v.Value)
	}

	if err := query.Select(); err != nil {
		return roomError(err)
	}

	return nil
}

func (r *roomRepository) GetByID(id int64, m interface{}) error {
	room, ok := m.(*model.Room)
	if !ok {
		return ErrInvalidType
	}
	room.ID = id

	if err := r.db.Conn().Model(room).WherePK().Select(); err != nil {
		return roomError(err)
	}

	return nil
}

func (r *roomRepository) GetBetween(start time.Time, end time.Time, m interface{}) error {
	return nil
}

func (r *roomRepository) DeleteByID(id int64) error {
	if _, err := r.db.Conn().Model(&model.Room{
		ID: id,
	}).WherePK().Delete(); err != nil {
		return err
	}
	return nil
}

func roomError(e error) error {
	pgErr, ok := e.(pg.Error)
	switch {
	case e == database.ErrorDNE:
		return ErrRoomDNE
	case ok && pgErr.IntegrityViolation():
		return ErrRoomExistsError
	default:
		return e
	}
}
