package service_test

import (
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/booking/config"
	"github.com/booking/logger"
	"github.com/booking/repository"
	"github.com/booking/service"

	"github.com/stretchr/testify/assert"

	"github.com/booking/model"
	"github.com/booking/repository/mocks"
)

func TestRoomService(t *testing.T) {
	c := &config.Config{}

	t.Run("Create", func(t *testing.T) {
		room := &model.Room{
			ID:      1,
			Name:    "C1",
			Company: "C",
			Number:  1,
		}

		r := &mocks.Repository{}
		r.On("Create", room).Return(nil)

		s := service.NewRoomService(c, r, logger.NewLogger(c).WithField("env", "test"))

		err := s.Create(room)

		assert.NoError(t, err)
		r.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("GetAll", func(t *testing.T) {
		expected := []model.Room{{
			ID:      1,
			Name:    "C1",
			Company: model.CompanyCoke,
			Number:  1,
		}}
		query := []repository.Query{{
			Model: "room",
			Field: "name",
			Value: "C1",
		}, {
			Model: "room",
			Field: "company",
			Value: model.CompanyCoke,
		}}

		r := &mocks.Repository{}
		r.On("Get", query, &[]model.Room{}).Run(func(a mock.Arguments) {
			rooms := a.Get(1).(*[]model.Room)
			(*rooms) = append(*rooms, expected[0])
		}).Return(nil)

		s := service.NewRoomService(c, r, logger.NewLogger(c).WithField("env", "test"))

		rooms, err := s.GetAll("C1", "coke")

		assert.NoError(t, err)
		assert.Equal(t, expected, rooms)
		r.AssertNumberOfCalls(t, "Get", 1)
	})

	t.Run("Get", func(t *testing.T) {
		id := int64(1)
		expected := &model.Room{
			ID:      id,
			Name:    "C1",
			Company: model.CompanyCoke,
			Number:  1,
		}

		r := &mocks.Repository{}
		r.On("GetByID", id, &model.Room{}).Run(func(a mock.Arguments) {
			room := a.Get(1).(*model.Room)
			(*room) = (*expected)
		}).Return(nil)

		s := service.NewRoomService(c, r, logger.NewLogger(c).WithField("env", "test"))

		room, err := s.Get(id)

		assert.NoError(t, err)
		assert.Equal(t, expected, room)
		r.AssertNumberOfCalls(t, "GetByID", 1)
	})

	t.Run("Delete", func(t *testing.T) {
		id := int64(1)

		r := &mocks.Repository{}
		r.On("DeleteByID", id).Return(nil)

		s := service.NewRoomService(c, r, logger.NewLogger(c).WithField("env", "test"))

		err := s.Delete(id)

		assert.NoError(t, err)
		r.AssertNumberOfCalls(t, "DeleteByID", 1)
	})
}
