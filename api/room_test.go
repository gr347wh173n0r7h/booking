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

	"github.com/booking/api"
	"github.com/booking/config"
	"github.com/booking/logger"
	"github.com/booking/model"
	"github.com/booking/service/mocks"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
)

var headers = http.Header{"Content-Type": []string{"application/json"}}

func TestAddRoomMissingErrors(t *testing.T) {
	u, _ := url.Parse("/rooms")

	svc := &mocks.RoomService{}
	a := api.NewRoomAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("MissingRoomNumber", func(t *testing.T) {
		j, err := json.Marshal(&model.RoomRequest{
			Company: "coke",
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
			errors.New("room number empty").Error(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"errors":%s}`, expected), rec.Body.String())
	})
	t.Run("InvalidCompanyName", func(t *testing.T) {
		j, err := json.Marshal(&model.RoomRequest{
			Number:  1,
			Company: "foo",
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
			errors.New("invalid company name").Error(),
		})
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, fmt.Sprintf(`{"errors":%s}`, expected), rec.Body.String())
	})
}

func TestAddRoom(t *testing.T) {
	svc := &mocks.RoomService{}
	a := api.NewRoomAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("AddRoom", func(t *testing.T) {
		rr := &model.RoomRequest{
			Number:  1,
			Company: "coke",
		}
		j, err := json.Marshal(rr)
		assert.NoError(t, err)

		svc.On("Create", rr.Model()).Return(nil)

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

func TestGetRooms(t *testing.T) {
	u, _ := url.Parse("/rooms/all?name=C1&company=C")

	rooms := []model.Room{{
		ID:      1,
		Name:    "C1",
		Number:  1,
		Company: model.CompanyCoke,
	}}

	svc := &mocks.RoomService{}
	a := api.NewRoomAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("GetRooms", func(t *testing.T) {
		svc.On("GetAll", "C1", "C").Return(rooms, nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "GET",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		expectedResponse, err := json.Marshal(rooms)
		assert.NoError(t, err)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponse), rec.Body.String())

	})
}

func TestGetRoom(t *testing.T) {
	u, _ := url.Parse("/rooms/1")

	room := model.Room{
		ID:      1,
		Name:    "C1",
		Number:  1,
		Company: model.CompanyCoke,
	}

	svc := &mocks.RoomService{}
	a := api.NewRoomAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

	c := restful.NewContainer()
	c.Add(a.WebService())

	t.Run("GetRoom", func(t *testing.T) {

		svc.On("Get", int64(1)).Return(&room, nil)

		rec := httptest.NewRecorder()
		req := restful.NewRequest(&http.Request{
			Header: headers,
			Method: "GET",
			URL:    u,
		})

		c.ServeHTTP(rec, req.Request)

		expectedResponse, err := json.Marshal(room)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, string(expectedResponse), rec.Body.String())

	})
}

func TestDeleteRoom(t *testing.T) {
	u, _ := url.Parse("/rooms/1")

	svc := &mocks.RoomService{}
	a := api.NewRoomAPI(svc, logger.NewLogger(&config.Config{}).WithField("env", "test"))

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
