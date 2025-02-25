package handler

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10" // For validation
	"github.com/segmentio/kafka-go"          // Kafka client library
	"str/internal"                           // Internal package
	"time"
)

// Metadata represents metadata information
type Metadata struct {
	Platform        string    `json:"platform" example:"booking.com"`                // Platform name
	SubmissionDate  time.Time `json:"submissionDate" example:"2024-07-21T17:32:28Z"` // Date of submission
	AdditionalProps struct{}  `json:"additionalProp1"`                               // Placeholder for additional properties
}

// Address represents an address with validation
type Address struct {
	Street     string `json:"street" example:"Culliganlaan 5"`                   // Street address
	City       string `json:"city" example:"Diegem"`                             // City name
	PostalCode string `json:"postalCode" example:"1831"`                         // Postal code
	Country    string `json:"country" example:"BEL" validate:"iso3166_1_alpha3"` // Country code with ISO 3166-1 alpha-3 validation
}

// Created represents a creation response
type Created struct {
	Status string `json:"status" example:"delivered"` // Status of the creation
}

// BadRequestError represents a bad request error
type BadRequestError struct {
	Status string `json:"status" example:"Wrong data format!"` // Error message for bad requests
}

// Error for validation represents an error for a validation check
type Errorvalidation struct {
	Type            string `json:"type" example:"string"`            // Type of the error
	Title           string `json:"title" example:"string"`           // Title of the error
	Status          int    `json:"status" example:"0"`               // Status code of the error
	Detail          string `json:"detail" example:"string"`          // Detailed error message
	Instance        string `json:"instance" example:"string"`        // Instance of the error
	AdditionalProp1 string `json:"additionalProp1" example:"string"` // Additional property 1
	AdditionalProp2 string `json:"additionalProp2" example:"string"` // Additional property 2
	AdditionalProp3 string `json:"additionalProp3" example:"string"` // Additional property 3
}

// UnauthorizedError represents an unauthorized error
type UnauthorizedError struct {
	Message string `json:"message" example:"JWT is invalid"` // Error message for unauthorized access
}

// numbervalid represents a valid number check
type  Numbervalid struct {
	Message string `json:"isValid" example:"True"` // Error message for a valid registration number
}

// InternalServerError represents an internal server error
type InternalServerError struct {
	Error string `json:"error" example:"An unexpected error occurred"` // Error message for internal server errors
}

// HTTPError represents a generic HTTP error
type HTTPError struct {
	Code    int    `json:"code"`    // HTTP status code
	Message string `json:"message"` // Error message
}

// NotFoundError represents a not found error
type NotFoundError struct {
	Error string `json:"error" example:"Resource not found"` // Error message for resource not found
}

// DeleteResponse represents a delete response
type DeleteResponse struct {
	Status string `json:"status"` // Status of the delete operation
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"` // Error message
}

// validationError represents validation error details
type validationError struct {
	Namespace       string `json:"namespace"`       // Namespace of the field
	Field           string `json:"field"`           // Field name
	StructNamespace string `json:"structNamespace"` // Struct namespace
	StructField     string `json:"structField"`     // Struct field name
	Tag             string `json:"tag"`             // Validation tag
	ActualTag       string `json:"actualTag"`       // Actual validation tag
	Kind            string `json:"kind"`            // Kind of field
	Type            string `json:"type"`            // Type of field
	Value           string `json:"value"`           // Value of the field
	Message         string `json:"message"`         // Error message
}

// Version of the application
var Version = "development"

var LOG_LEVEL = "DEBUG"

// Single instance of Validate for caching struct info
var validate *validator.Validate

// Status represents a status response
type Status struct {
	Status string `json:"status" example:"ok"` // Status message
}

// Identity represents identity information
type Identity struct {
	OAuth2AppName string   `json:"oauth2_app_name"` // OAuth2 application name
	CA            []string `json:"ca"`              // Certificate authorities
}

// ValidateStruct validates any struct based on tags set
func ValidateStruct(v *validator.Validate, s interface{}) (string, error) {
	err := v.Struct(s) // Validates the struct
	if err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			internal.Fatalf(err.Error()) // Logs fatal error if validation is invalid
		}

		for _, err := range err.(validator.ValidationErrors) { // Iterates over validation errors
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
				Message:         err.Error(), // Sets the error message
			}

			return e.Message, nil // Returns the first validation error message
		}
	}
	return "", nil // Returns no error if validation passes
}

// convertHeadersToBytes converts Kafka headers to a byte slice
func convertHeadersToBytes(headers []kafka.Header) []byte {
	var buffer bytes.Buffer
	for _, header := range headers {
		buffer.Write(header.Value) // Writes each header value to the buffer
	}
	return buffer.Bytes() // Returns the byte slice
}

// Error status codes
const (
    ErrBadRequest          = "BAD_REQUEST"
    ErrUnauthorized        = "UNAUTHORIZED"
    ErrValidation          = "VALIDATION"
    ErrServiceUnavailable  = "SERVICE_UNAVAILABLE"
    ErrNotFound           = "NOT_FOUND"
)

// Error IDs for specific scenarios
const (
    // Bad Request Errors (400)
    ErrIDInvalidJSON          = "BAD_REQUEST_001"
    ErrIDInvalidParam         = "BAD_REQUEST_002"
    ErrIDMissingField         = "BAD_REQUEST_003"
    
    // Auth Errors (401)
    ErrIDAuthFailed           = "AUTH_001"
    
    // Validation Errors (422)
    ErrIDValidationFailed     = "VALIDATION_001"
    ErrIDDataValidation       = "VALIDATION_002"
    ErrIDStructValidation     = "VALIDATION_003"
    ErrIDFormatValidation     = "VALIDATION_004"
    
    // Service Errors (503)
    ErrIDServiceFailed        = "SERVICE_001"
    ErrIDExternalAPIFailed    = "SERVICE_002"
    ErrIDKafkaWriteFailed     = "SERVICE_003"
    
    // Not Found Errors (404)
    ErrIDResourceNotFound     = "NOT_FOUND_001"
)

// APIError represents an API error response
type APIError struct {
    Status  string       `json:"status"`
    Results []ErrorEntry `json:"results"`
}

// ErrorEntry represents a detailed error entry
type ErrorEntry struct {
    Message      string `json:"message"`
    Code         int    `json:"code"`
    ErrorUserMsg string `json:"error_user_msg"`
    ID           string `json:"id"`
}

// NewAPIError creates a new APIError instance
func NewAPIError(message string, code int, userMsg string, id string) APIError {
    return APIError{
        Status: "error",
        Results: []ErrorEntry{
            {
                Message:      message,
                Code:         code,
                ErrorUserMsg: userMsg,
                ID:          id,
            },
        },
    }
}