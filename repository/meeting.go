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
	// ErrMeetingDNE defined a Meeting does not exist error
	ErrMeetingDNE error = errors.New("meeting does not exist")
	// ErrMeetingExistsError defined a Meeting already exists
	ErrMeetingExistsError error = errors.New("meeting already exist")
)

type meetingRepository struct {
	db database.Database
}

// NewMeetingRepository returns a meeting implementation of Repository
func NewMeetingRepository(db database.Database, log bool) (Repository, error) {
	if err := db.CreateSchema([]interface{}{
		(*model.Meeting)(nil),
	}); err != nil {
		return nil, err
	}

	if log {
		db.Conn().AddQueryHook(dbLogger{})
	}

	return &meetingRepository{
		db: db,
	}, nil
}

func (r *meetingRepository) Create(m interface{}) error {
	meeting, ok := m.(*model.Meeting)
	if !ok {
		return ErrInvalidType
	}

	exists, err := r.db.Conn().Model(&[]model.Meeting{}).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.Where("meeting.room_id = ?", meeting.RoomID).
				Where("meeting.start <= ?", meeting.Start).
				Where("meeting.end >= ?", meeting.Start)
			return q, nil
		}).
		WhereOrGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.Where("meeting.room_id = ?", meeting.RoomID).
				Where("meeting.start <= ?", meeting.End).
				Where("meeting.end >= ?", meeting.End)
			return q, nil
		}).Exists()

	if err != nil {
		return err
	}

	if exists {
		return ErrMeetingExistsError
	}

	_, err = r.db.Conn().Model(meeting).Insert()
	return meetingError(err)
}

func (r *meetingRepository) Get(q []Query, m interface{}) error {
	meetings, ok := m.(*[]model.Meeting)
	if !ok {
		return ErrInvalidType
	}

	query := r.db.Conn().Model(meetings)

	for _, v := range q {
		query = query.Where(fmt.Sprintf("%v.%v = ?", v.Model, v.Field), v.Value)
	}

	if err := query.Relation("Room").Select(); err != nil {
		return meetingError(err)
	}

	return nil
}

func (r *meetingRepository) GetByID(id int64, m interface{}) error {
	meeting, ok := m.(*model.Meeting)
	if !ok {
		return ErrInvalidType
	}
	meeting.ID = id

	if err := r.db.Conn().Model(meeting).Relation("Room").WherePK().Select(); err != nil {
		return meetingError(err)
	}

	return nil
}

func (r *meetingRepository) GetBetween(start time.Time, end time.Time, m interface{}) error {
	meetings, ok := m.(*[]model.Meeting)
	if !ok {
		return ErrInvalidType
	}

	query := r.db.Conn().Model(meetings).
		WhereGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.Where("meeting.start >= ?", start).
				Where("meeting.start <= ?", end)
			return q, nil
		}).
		WhereOrGroup(func(q *pg.Query) (*pg.Query, error) {
			q = q.Where("meeting.end >= ?", start).
				Where("meeting.end <= ?", end)
			return q, nil
		})

	if err := query.Select(); err != nil {
		return meetingError(err)
	}

	return nil
}

func (r *meetingRepository) DeleteByID(id int64) error {
	if _, err := r.db.Conn().Model(&model.Meeting{
		ID: id,
	}).WherePK().Delete(); err != nil {
		return err
	}
	return nil
}

func meetingError(e error) error {
	pgErr, ok := e.(pg.Error)
	switch {
	case e == database.ErrorDNE:
		return ErrMeetingDNE
	case ok && pgErr.IntegrityViolation():
		switch pgErr.Field('C') {
		case "23503":
			return ErrRoomDNE
		default:
			return ErrMeetingExistsError
		}
	default:
		return e
	}
}
