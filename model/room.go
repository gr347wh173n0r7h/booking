package model

import (
	"errors"
	"fmt"
	"strings"
)

// Company defines a company enum
type Company string

const (
	// CompanyCoke defines Coke as a Company
	CompanyCoke Company = "C"
	// CompanyPepsi defines Pepsi as a Company
	CompanyPepsi Company = "P"
)

var (
	// CompanyName represents a map to convert Company to string
	CompanyName = map[Company]string{
		"C": "coke",
		"P": "pepsi",
	}
	// CompanyID represents a map to convert string to Company
	CompanyID = map[string]Company{
		"coke":  "C",
		"pepsi": "P",
	}
)

// ModelRoom defines Room model name for go-pg
const ModelRoom = "room"

// Room defines a storable meeting structure
type Room struct {
	ID      int64
	Name    string
	Number  int     `pg:",unique:vector"`
	Company Company `pg:",unique:vector"`
}

func (r Room) String() string {
	return fmt.Sprintf("Room<%d %s %v>", r.ID, r.Name, r.Number)
}

// RoomRequest defines expected Room request
type RoomRequest struct {
	Number  int
	Company string
}

// Validate validates contents of RoomRequest
func (r *RoomRequest) Validate() error {
	if r.Number == 0 {
		return errors.New("room number empty")
	}
	if r.Company == "" {
		return errors.New("Company empty")
	}
	if _, ok := CompanyID[strings.ToLower(r.Company)]; !ok {
		return errors.New("invalid company name")
	}
	return nil
}

// Model transforms RoomRequest to Room
func (r *RoomRequest) Model() *Room {
	cid := CompanyID[strings.ToLower(r.Company)]
	return &Room{
		Name:    fmt.Sprintf("%s%d", string(cid), r.Number),
		Number:  r.Number,
		Company: cid,
	}
}
