package service

import (
	"github.com/sirupsen/logrus"

	"github.com/booking/config"
	"github.com/booking/model"
	"github.com/booking/repository"
)

// RoomService defines interface for services creating Rooms
type RoomService interface {
	Create(r *model.Room) error
	GetAll(roomName string, companyName string) ([]model.Room, error)
	Get(id int64) (*model.Room, error)
	Delete(id int64) error
}

type roomService struct {
	config *config.Config
	repo   repository.Repository
	logger *logrus.Entry
}

// NewRoomService returns a roomService implementation of RoomService
func NewRoomService(c *config.Config, r repository.Repository, l *logrus.Entry) RoomService {
	return &roomService{
		config: c,
		repo:   r,
		logger: l,
	}
}

func (s *roomService) Create(r *model.Room) error {
	return s.repo.Create(r)
}

func (s *roomService) GetAll(roomName string, companyName string) ([]model.Room, error) {
	query := []repository.Query{}
	if roomName != "" {
		query = append(query, repository.Query{
			Model: model.ModelRoom,
			Field: "name",
			Value: roomName,
		})
	}
	if companyName != "" {
		query = append(query, repository.Query{
			Model: model.ModelRoom,
			Field: "company",
			Value: model.CompanyID[companyName],
		})
	}

	rooms := []model.Room{}
	err := s.repo.Get(query, &rooms)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (s *roomService) Get(id int64) (*model.Room, error) {
	room := &model.Room{}
	err := s.repo.GetByID(id, room)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (s *roomService) Delete(id int64) error {
	return s.repo.DeleteByID(id)
}
