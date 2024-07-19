package handler

import (
	"bytes"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"str/internal"
	"time"
)

type Metadata struct {
	Platform       string    `json:"platform" example:"booking.com"`
	SubmissionDate time.Time `json:"submissionDate" example:"2024-07-21T17:32:28Z"`
	AdditionalProp struct{}  `json:"additionalProp1"`
}

type Address struct {
	Street     string `json:"street" example:"Culliganlaan 5"`
	City       string `json:"city" example:"Diegem"`
	PostalCode string `json:"postalCode" example:"1831"`
	Country    string `json:"country" example:"BEL" validate:"iso3166_1_alpha3"`
}

type validationError struct {
	Namespace       string `json:"namespace"` // can differ when a custom TagNameFunc is registered or
	Field           string `json:"field"`     // by passing alt name to ReportError like below
	StructNamespace string `json:"structNamespace"`
	StructField     string `json:"structField"`
	Tag             string `json:"tag"`
	ActualTag       string `json:"actualTag"`
	Kind            string `json:"kind"`
	Type            string `json:"type"`
	Value           string `json:"value"`
	Param           string `json:"param"`
	Message         string `json:"message"`
}

var Version = "development"

// use a single instance of Validate, it caches struct info
var validate *validator.Validate

type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

type Status struct {
	Status string `json:"status" example:"ok"`
}

type Identity struct {
	OAuth2AppName string `json:"oauth2_app_name"`
	CA            string `json:"ca"`
}

// ValidateStruct validates any struct passed to it based on tags set
func ValidateStruct(v *validator.Validate, s interface{}) (string, error) {
	err := v.Struct(s)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			internal.Fatalf(err.Error())
		}

		for _, err := range err.(validator.ValidationErrors) {
			e := validationError{
				Namespace:       err.Namespace(),
				Field:           err.Field(),
				StructNamespace: err.StructNamespace(),
				StructField:     err.StructField(),
				Tag:             err.Tag(),
				ActualTag:       err.ActualTag(),
				Kind:            fmt.Sprintf("%v", err.Kind()),
				Type:            fmt.Sprintf("%v", err.Type()),
				Value:           fmt.Sprintf("%v", err.Value()),
				Param:           err.Param(),
				Message:         err.Error(),
			}

			return e.Message, nil
		}
	}
	return "", nil
}

func convertHeadersToBytes(headers []kafka.Header) []byte {
	var buffer bytes.Buffer
	for _, header := range headers {
		buffer.Write(header.Value)
	}
	return buffer.Bytes()
}
