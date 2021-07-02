package api

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

// API defines interfaces to provide webservices
type API interface {
	WebService() *restful.WebService
}

// EnrichSwaggerObject enriches swagger object with API metadata
func EnrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "BookingService",
			Description: "Resource for creating and booking meeting rooms",
			Contact: &spec.ContactInfo{
				Name:  "Jordan Petersen",
				Email: "jordan.a.petersen@gmail.com",
				URL:   "http://github.com/gr347wh173n0r7h",
			},
			License: &spec.License{
				Name: "MIT",
				URL:  "http://mit.org",
			},
			Version: "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{
		{TagProps: spec.TagProps{
			Name:        "Rooms",
			Description: "Managing meeting rooms",
		}},
		{TagProps: spec.TagProps{
			Name:        "Booking",
			Description: "Managing booking meetings",
		}},
	}
}
