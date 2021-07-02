package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/booking/api"
	"github.com/booking/config"
	"github.com/booking/logger"
	"github.com/booking/model"
	"github.com/booking/service/mocks"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
)

func TestAddMeetingMissingErrors(t *testing.T) {
	u, _ := url.Parse("/booking/")

	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("MissingRoomID", func(t *testing.T) {
		j, err := json.Marshal(&model.Meeting{})
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "POST",
			URL:    u,
			Body:   ioutil.NopCloser(bytes.NewReader(j)),
		})

		c.ServeHTTP(rec, req.Request)

		expected, err := json.Marshal([]string{
			errors.New("room-id empty").Error(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"errors":%s}`, expected), rec.Body.String())
	})
	t.Run("InvalidTitle", func(t *testing.T) {
		j, err := json.Marshal(&model.Meeting{
			RoomID: 1,
		})
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "POST",
			URL:    u,
			Body:   ioutil.NopCloser(bytes.NewReader(j)),
		})

		c.ServeHTTP(rec, req.Request)

		expected, err := json.Marshal([]string{
			errors.New("title empty").Error(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"errors":%s}`, expected), rec.Body.String())
	})
}

func TestAddMeeting(t *testing.T) {
	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("AddMeeting", func(t *testing.T) {
		start := time.Now()
		mr := &model.MeetingRequest{
			RoomID: 1,
			Title:  "foo",
			Start:  &start,
		}
		j, err := json.Marshal(mr)

		svc.On("Create", mr.Model()).Return(nil)

		rec := httptest.NewRecorder()
		res := restful.NewResponse(rec)
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "POST",
			URL:    &url.URL{},
			Body:   ioutil.NopCloser(bytes.NewReader(j)),
		})

		c.ServeHTTP(rec, req.Request)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, res.StatusCode())
	})
}

func TestGetMeetings(t *testing.T) {
	u, _ := url.Parse("/booking/meetings/all?room-id=1")

	meetings := []model.Meeting{{
		ID:     1,
		RoomID: 1,
	}}

	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("GetMeetings", func(t *testing.T) {
		svc.On("GetAll", 1).Return(meetings, nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "GET",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		expectedResponse, err := json.Marshal(meetings)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponse), rec.Body.String())

	})
}

func TestGetMeeting(t *testing.T) {
	u, _ := url.Parse("/booking/meetings/1")

	meeting := model.Meeting{
		ID:     1,
		RoomID: 1,
	}

	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("GetMeeting", func(t *testing.T) {

		svc.On("Get", int64(1)).Return(&meeting, nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "GET",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		expectedResponse, err := json.Marshal(meeting)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponse), rec.Body.String())

	})
}

func TestDeleteMeetings(t *testing.T) {
	u, _ := url.Parse("/booking/meetings/1")

	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("DeleteMeeting", func(t *testing.T) {

		svc.On("Delete", int64(1)).Return(nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "DELETE",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		assert.Equal(t, http.StatusOK, rec.Code)

	})
}

func TestGetAvailable(t *testing.T) {
	u, _ := url.Parse("/booking/available?date=2021-07-01T02%3A43%3A21%2B00%3A00")
	date, _ := time.Parse(time.RFC3339, "2021-07-01T02:43:21+00:00")

	am := model.AvailabilityMap{
		1: map[time.Time]*model.Meeting{
			time.Now(): {ID: 1},
		},
	}

	svc := &mocks.BookingService{}
	a := api.NewBookingAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("GetAvailable", func(t *testing.T) {

		svc.On("GetAvailable", date.UTC()).Return(am, nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "GET",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		assert.Equal(t, http.StatusOK, rec.Code)
		expectedResponse, err := json.Marshal(am)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponse), rec.Body.String())
	})
}
