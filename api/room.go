package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/sirupsen/logrus"

	"github.com/booking/model"

	"github.com/booking/service"
)

// RoomRootPath represents base room path
const RoomRootPath = "/rooms"

var roomTags = []string{"Rooms"}

type roomAPI struct {
	service service.RoomService
	logger  *logrus.Entry
}

// NewRoomAPI returns a roomAPI implementation of API
func NewRoomAPI(s service.RoomService, l *logrus.Entry) API {
	return &roomAPI{
		service: s,
		logger:  l,
	}
}

func (a *roomAPI) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(RoomRootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(
		ws.POST("/").To(a.AddRoomHandler).
			Doc("add room").
			Metadata(restfulspec.KeyOpenAPITags, roomTags).
			Reads(model.RoomRequest{}).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.Room{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), []error{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []error{}),
	)
	ws.Route(
		ws.GET("/all").To(a.GetRoomsHandler).
			Doc("get all rooms").
			Metadata(restfulspec.KeyOpenAPITags, roomTags).
			Param(ws.QueryParameter("name", "Room name").
				DataType("string").
				Required(false).
				AllowMultiple(false)).
			Param(ws.QueryParameter("company", "Company name").
				DataType("string").
				Required(false).
				AllowMultiple(false)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), []model.Room{}),
	)
	ws.Route(
		ws.GET("/{room-id}").To(a.GetRoomHandler).
			Doc("get room by id").
			Metadata(restfulspec.KeyOpenAPITags, roomTags).
			Param(ws.PathParameter("room-id", "identifier of room").
				DataType("string")).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.Room{}),
	)
	ws.Route(
		ws.DELETE("/{room-id}").To(a.DeleteRoomHandler).
			Doc("delete room by id").
			Metadata(restfulspec.KeyOpenAPITags, roomTags).
			Param(ws.PathParameter("room-id", "identifier of room").
				DataType("string")).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), nil),
	)

	return ws
}

func (a *roomAPI) AddRoomHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "AddRoomHandler").
		WithField("body", req.Request.Body)

	log.Debug("begin handler")
	defer log.Debug("end handler")

	room := &model.RoomRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(room); err != nil {
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err := room.Validate(); err != nil {
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err := a.service.Create(room.Model()); err != nil {
		log.WithError(err).Error("error adding rooms")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (a *roomAPI) GetRoomsHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "GetRoomsHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	name := req.QueryParameter("name")
	company := req.QueryParameter("company")

	rooms, err := a.service.GetAll(name, company)
	if err != nil {
		log.WithError(err).Error("error getting rooms")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	WriteJSON(res, a.logger, rooms)
}

func (a *roomAPI) GetRoomHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "GetRoomHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	roomID, err := strconv.Atoi(req.PathParameter("room-id"))
	if err != nil {
		log.WithError(err).Error("invalid room-id")
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	room, err := a.service.Get(int64(roomID))
	if err != nil {
		log.WithError(err).Error("error getting room")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	WriteJSON(res, a.logger, room)
}

func (a *roomAPI) DeleteRoomHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "DeleteRoomHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	roomID, err := strconv.Atoi(req.PathParameter("room-id"))
	if err != nil {
		a.logger.WithError(err).Error("invalid room-id")
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err = a.service.Delete(int64(roomID)); err != nil {
		log.WithError(err).Error("error deleting room")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}
