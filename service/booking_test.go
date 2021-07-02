package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/booking/config"
	"github.com/booking/logger"
	"github.com/booking/model"
	"github.com/booking/repository"
	"github.com/booking/repository/mocks"
	"github.com/booking/service"
)

func TestBookingService(t *testing.T) {
	c := &config.Config{MaxTimeBlockMin: 60}
	t.Run("Create", func(t *testing.T) {
		meeting := model.Meeting{
			ID:     1,
			RoomID: 2,
		}

		mr := &mocks.Repository{}
		mr.On("Create", &meeting).Return(nil)
		rr := &mocks.Repository{}

		s := service.NewBookingService(c, mr, rr, logger.NewLogger(c).WithField("env", "test"))

		err := s.Create(&meeting)

		assert.NoError(t, err)
		mr.AssertNumberOfCalls(t, "Create", 1)
	})

	t.Run("GetAll", func(t *testing.T) {
		expected := []model.Meeting{{
			ID:     1,
			RoomID: 1,
		}, {
			ID:     1,
			RoomID: 1,
		}}
		query := []repository.Query{{
			Model: "meeting",
			Field: "room_id",
			Value: 1,
		}}

		mr := &mocks.Repository{}
		mr.On("Get", query, &[]model.Meeting{}).Run(func(a mock.Arguments) {
			meetings := a.Get(1).(*[]model.Meeting)
			(*meetings) = append(*meetings, expected...)
		}).Return(nil)
		rr := &mocks.Repository{}

		s := service.NewBookingService(c, mr, rr, logger.NewLogger(c).WithField("env", "test"))

		meetings, err := s.GetAll(1)

		assert.NoError(t, err)
		assert.Equal(t, expected, meetings)
		mr.AssertNumberOfCalls(t, "Get", 1)
	})

	t.Run("Get", func(t *testing.T) {
		id := int64(1)
		expected := &model.Meeting{
			ID:     1,
			RoomID: 1,
		}

		mr := &mocks.Repository{}
		mr.On("GetByID", id, &model.Meeting{}).Run(func(a mock.Arguments) {
			meeting := a.Get(1).(*model.Meeting)
			(*meeting) = (*expected)
		}).Return(nil)
		rr := &mocks.Repository{}

		s := service.NewBookingService(c, mr, rr, logger.NewLogger(c).WithField("env", "test"))

		meeting, err := s.Get(id)

		assert.NoError(t, err)
		assert.Equal(t, expected, meeting)
		mr.AssertNumberOfCalls(t, "GetByID", 1)
	})

	t.Run("Delete", func(t *testing.T) {
		id := int64(1)

		mr := &mocks.Repository{}
		mr.On("DeleteByID", id).Return(nil)
		rr := &mocks.Repository{}

		s := service.NewBookingService(c, mr, rr, logger.NewLogger(c).WithField("env", "test"))

		err := s.Delete(id)

		assert.NoError(t, err)
		mr.AssertNumberOfCalls(t, "DeleteByID", 1)
	})

	t.Run("GetAvailable", func(t *testing.T) {
		sTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)
		m1Time := time.Date(0, 0, 0, 1, 0, 0, 0, time.UTC)
		m2Time := time.Date(0, 0, 0, 10, 0, 0, 0, time.UTC)

		availableRooms := []model.Room{{
			ID: 1,
		}, {
			ID: 2,
		}}

		currentMeetings := []model.Meeting{{
			ID:     3,
			RoomID: 1,
			Start:  m1Time,
		}, {
			ID:     4,
			RoomID: 2,
			Start:  m2Time,
		}, {
			ID:     5,
			RoomID: 2,
			Start:  m1Time,
		}}

		expected := model.AvailabilityMap{}
		for i, m := range currentMeetings {
			if _, ok := expected[m.RoomID]; !ok {
				expected[m.RoomID] = map[time.Time]*model.Meeting{}
				for _, tv := range model.CreateTimeSlotMap(sTime, c.MaxTimeBlockMin) {
					expected[m.RoomID][tv] = nil
				}
			}
			expected[m.RoomID][m.Start] = &currentMeetings[i]
		}

		rr := &mocks.Repository{}
		rr.On("Get", []repository.Query{}, &[]model.Room{}).Run(func(a mock.Arguments) {
			rooms := a.Get(1).(*[]model.Room)
			(*rooms) = append(*rooms, availableRooms...)
		}).Return(nil)

		mr := &mocks.Repository{}
		mr.On(
			"GetBetween",
			time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC),
			time.Date(0, 0, 0, 24, 0, 0, 0, time.UTC),
			&[]model.Meeting{},
		).Run(func(a mock.Arguments) {
			meetings := a.Get(2).(*[]model.Meeting)
			(*meetings) = append(*meetings, currentMeetings...)
		}).Return(nil)

		s := service.NewBookingService(c, mr, rr, logger.NewLogger(c).WithField("env", "test"))

		am, err := s.GetAvailable(sTime)

		assert.NoError(t, err)
		assert.Equal(t, expected, am)
		mr.AssertNumberOfCalls(t, "GetBetween", 1)
	})
}
