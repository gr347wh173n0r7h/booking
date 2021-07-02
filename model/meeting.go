package model

import (
	"errors"
	"fmt"
	"time"
)

// ModelMeeting defines Meeting model name for go-pg
const ModelMeeting = "meeting"

// Meeting defines a storable meeting structure
type Meeting struct {
	ID        int64
	RoomID    int64 `pg:"on_delete:CASCADE"`
	Room      *Room `pg:"rel:has-one"`
	Title     string
	Attendees []string
	Created   time.Time `pg:"default:now()"`
	Start     time.Time
	End       time.Time
}

func (m Meeting) String() string {
	return fmt.Sprintf("Meeting<%d %d %s>", m.ID, m.RoomID, m.Title)
}

// AvailabilityMap defines a map of available Room and Time slots with corresponding Meeting if booked
//  else Time slots will be nil
type AvailabilityMap map[int64]map[time.Time]*Meeting

// MeetingRequest defines a expected Meeting request
type MeetingRequest struct {
	RoomID    int64
	Title     string
	Attendees []string
	Start     *time.Time
}

// Validate validates contents of MeetingRequest
func (r *MeetingRequest) Validate() error {
	if r.RoomID == 0 {
		return errors.New("room-id empty")
	}
	if r.Title == "" {
		return errors.New("title empty")
	}
	if r.Start == nil {
		return errors.New("start empty")
	}
	if *r.Start != time.Date(r.Start.Year(), r.Start.Month(), r.Start.Day(), r.Start.Hour(), 0, 0, 0, r.Start.Location()) {
		return errors.New("invalid start time")
	}
	return nil
}

// Model transforms MeetingRequest to Meeting
func (r *MeetingRequest) Model() *Meeting {
	return &Meeting{
		RoomID:    r.RoomID,
		Title:     r.Title,
		Attendees: r.Attendees,
		Start:     *r.Start,
	}
}

// CreateTimeSlotMap creates a slice of time blocks for requested interval
func CreateTimeSlotMap(date time.Time, maxTimeBlock int) []time.Time {
	ts := []time.Time{}
	t := time.Date(date.Year(), date.Month(), date.Day(),
		0, 0, 0, 0, time.UTC)
	for i := 0; i < (24 * (60 / maxTimeBlock)); i++ {
		ts = append(ts, t)
		t = t.Add(time.Minute * time.Duration(maxTimeBlock))
	}
	return ts
}
