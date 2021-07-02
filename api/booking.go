package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/sirupsen/logrus"

	"github.com/booking/model"

	"github.com/booking/service"
)

// BookingRootPath represents base booking path
const BookingRootPath = "/booking"

type bookingAPI struct {
	service service.BookingService
	logger  *logrus.Entry
}

var bookingTags = []string{"Booking"}

// NewBookingAPI returns a bookingAPI implementation of API
func NewBookingAPI(s service.BookingService, l *logrus.Entry) API {
	return &bookingAPI{
		service: s,
		logger:  l,
	}
}

func (a *bookingAPI) WebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.Path(BookingRootPath).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(
		ws.POST("/").To(a.AddMeetingHandler).
			Doc("add meeting").
			Metadata(restfulspec.KeyOpenAPITags, bookingTags).
			Reads(model.MeetingRequest{}).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.Meeting{}).
			Returns(http.StatusBadRequest, http.StatusText(http.StatusBadRequest), []error{}).
			Returns(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), []error{}),
	)
	ws.Route(
		ws.GET("/meetings/all").To(a.GetMeetingsHandler).
			Doc("get all meetings").
			Metadata(restfulspec.KeyOpenAPITags, bookingTags).
			Param(ws.QueryParameter("room-id", "Room name").
				DataType("string").
				Required(false).
				AllowMultiple(false)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), []model.Meeting{}),
	)
	ws.Route(
		ws.GET("/meetings/{meeting-id}").To(a.GetMeetingHandler).
			Doc("get meeting by id").
			Metadata(restfulspec.KeyOpenAPITags, bookingTags).
			Param(ws.PathParameter("meeting-id", "identifier of meeting").
				DataType("string")).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.Meeting{}),
	)
	ws.Route(
		ws.DELETE("/meetings/{meeting-id}").To(a.DeleteMeetingHandler).
			Doc("delete meeting by id").
			Metadata(restfulspec.KeyOpenAPITags, bookingTags).
			Param(ws.PathParameter("meeting-id", "identifier of meeting").
				DataType("string")).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.Meeting{}),
	)
	ws.Route(
		ws.GET("/available").To(a.GetAvailableHandler).
			Doc("get all availability").
			Metadata(restfulspec.KeyOpenAPITags, bookingTags).
			Param(ws.QueryParameter("date", "date for availability").
				DataType("string").
				Required(false).
				AllowMultiple(false)).
			Returns(http.StatusOK, http.StatusText(http.StatusOK), model.AvailabilityMap{}),
	)

	return ws
}

func (a *bookingAPI) AddMeetingHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "AddBookingHandler").
		WithField("body", req.Request.Body)

	log.Debug("begin handler")
	defer log.Debug("end handler")

	meeting := &model.MeetingRequest{}
	if err := json.NewDecoder(req.Request.Body).Decode(meeting); err != nil {
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err := meeting.Validate(); err != nil {
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err := a.service.Create(meeting.Model()); err != nil {
		log.WithError(err).Error("error adding meeting")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (a *bookingAPI) GetMeetingsHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "GetMeetingsHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	var roomID int
	roomQuery := req.QueryParameter("room-id")
	if roomQuery != "" {
		roomID, _ = strconv.Atoi(req.QueryParameter("room-id"))
	}

	meetings, err := a.service.GetAll(roomID)
	if err != nil {
		log.WithError(err).Error("error getting meetings")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	WriteJSON(res, a.logger, meetings)
}

func (a *bookingAPI) GetMeetingHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "GetMeetingHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	meetingID, err := strconv.Atoi(req.PathParameter("meeting-id"))
	if err != nil {
		log.WithError(err).Error("invalid meeting-id")
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	meeting, err := a.service.Get(int64(meetingID))
	if err != nil {
		log.WithError(err).Error("error getting meeting")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	WriteJSON(res, a.logger, meeting)
}

func (a *bookingAPI) DeleteMeetingHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "DeleteMeetingHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	meetingID, err := strconv.Atoi(req.PathParameter("meeting-id"))
	if err != nil {
		a.logger.WithError(err).Error("invalid meeting-id")
		WriteError(res, http.StatusBadRequest, a.logger, err)
		return
	}

	if err = a.service.Delete(int64(meetingID)); err != nil {
		log.WithError(err).Error("error deleting meeting")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	res.WriteHeader(http.StatusOK)
}

func (a *bookingAPI) GetAvailableHandler(req *restful.Request, res *restful.Response) {
	log := a.logger.WithField("handler", "GetAvailableHandler").
		WithField("params", req.PathParameters())

	log.Debug("begin handler")
	defer log.Debug("end handler")

	date := time.Now()
	dateStr := req.QueryParameter("date")
	if dateStr != "" {
		date, _ = time.Parse(time.RFC3339, dateStr)
	}

	meetings, err := a.service.GetAvailable(date.UTC())
	if err != nil {
		log.WithError(err).Error("error getting meetings")
		WriteError(res, http.StatusInternalServerError, a.logger, err)
		return
	}
	WriteJSON(res, a.logger, meetings)
}
