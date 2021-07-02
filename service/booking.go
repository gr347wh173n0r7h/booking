package service

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/booking/config"
	"github.com/booking/model"
	"github.com/booking/repository"
)

// BookingService defines interface for services booking Rooms for Meetings
type BookingService interface {
	Create(r *model.Meeting) error
	GetAll(roomID int) ([]model.Meeting, error)
	Get(id int64) (*model.Meeting, error)
	Delete(id int64) error
	GetAvailable(date time.Time) (model.AvailabilityMap, error)
}

type bookingService struct {
	config      *config.Config
	meetingRepo repository.Repository
	roomRepo    repository.Repository
	logger      *logrus.Entry
}

// NewBookingService returns a bookingService implementation of BookingService
func NewBookingService(c *config.Config, meetingRepo repository.Repository, roomRepo repository.Repository, l *logrus.Entry) BookingService {
	return &bookingService{
		config:      c,
		meetingRepo: meetingRepo,
		roomRepo:    roomRepo,
		logger:      l,
	}
}

func (s *bookingService) Create(r *model.Meeting) error {
	r.End = time.Date(r.Start.Year(), r.Start.Month(), r.Start.Day(), r.Start.Hour(), s.config.MaxTimeBlockMin, 0, 0, r.Start.Location())
	return s.meetingRepo.Create(r)
}

func (s *bookingService) GetAll(roomID int) ([]model.Meeting, error) {
	query := []repository.Query{}
	if roomID != 0 {
		query = append(query, repository.Query{
			Model: model.ModelMeeting,
			Field: "room_id",
			Value: roomID,
		})
	}

	meetings := []model.Meeting{}
	err := s.meetingRepo.Get(query, &meetings)
	if err != nil {
		return nil, err
	}
	return meetings, nil
}

func (s *bookingService) Get(id int64) (*model.Meeting, error) {
	meeting := &model.Meeting{}
	err := s.meetingRepo.GetByID(id, meeting)
	if err != nil {
		return nil, err
	}
	return meeting, nil
}

func (s *bookingService) Delete(id int64) error {
	return s.meetingRepo.DeleteByID(id)
}

func (s *bookingService) GetAvailable(date time.Time) (model.AvailabilityMap, error) {
	am := model.AvailabilityMap{}

	// Get all rooms
	rooms := []model.Room{}
	if err := s.roomRepo.Get([]repository.Query{}, &rooms); err != nil {
		return nil, err
	}

	// Create time slots
	ts := model.CreateTimeSlotMap(date, s.config.MaxTimeBlockMin)

	// Create room and time slots to availability map
	for _, r := range rooms {
		am[r.ID] = map[time.Time]*model.Meeting{}
		for _, t := range ts {
			am[r.ID][t] = nil
		}
	}

	// Get all meetings on date
	meetings := []model.Meeting{}
	if err := s.meetingRepo.GetBetween(
		time.Date(date.Year(), date.Month(), date.Day(),
			0, 0, 0, 0, time.UTC),
		time.Date(date.Year(), date.Month(), date.Day(),
			24, 0, 0, 0, time.UTC),
		&meetings,
	); err != nil {
		return nil, err
	}

	// Loop over meeting and remove timeslots from availability map.
	for i, m := range meetings {
		mt := time.Date(m.Start.Year(), m.Start.Month(), m.Start.Day(),
			m.Start.Hour(), 0, 0, 0, time.UTC)
		if _, ok := am[m.RoomID][mt]; ok {
			am[m.RoomID][mt] = &meetings[i]
		}
	}

	return am, nil
}
