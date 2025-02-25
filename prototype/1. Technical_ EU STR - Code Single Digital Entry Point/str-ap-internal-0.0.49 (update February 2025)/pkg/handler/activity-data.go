package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"io"
	"log/slog"
	"net/http"
	"os"
	"str/internal"
	"str/pkg/middleware"
	"strconv"
	"strings"
	"time"
)

type Temporal struct {
	StartDateTime time.Time `json:"startDateTime" example:"2024-07-21T17:32:28Z"`
	EndDateTime   time.Time `json:"endDateTime" example:"2024-07-25T17:32:28Z"`
}

type GuestData struct {
	Address                     Address  `json:"address"`
	CompetentAuthorityId_area   string   `json:"competentAuthorityId_area" example:"competentAuthorityId_area"`
	CompetentAuthorityName_area string   `json:"competentAuthorityName_area" example:"competentAuthorityName_area"`
	CountryOfGuests             []string `json:"countryOfGuests" example:"ITA,NLD"`
	HostID                      string   `json:"hostId" example:"placeholder-host-id"`
	NumberOfGuests              uint     `json:"numberOfGuests" example:"3"`
	Temporal                    Temporal `json:"temporal"`
	URL                         string   `json:"URL" example:"placeholder-URL"`
	RegistrationNumber          string   `json:"registrationNumber" example:"placeholder-registrationNumber"`
}

type CompetentAuthority struct {
	Code string `json:"competentAuthorityCode"`
	Name string `json:"competentAuthorityName"`
}

type ActivityData struct {
	Data     []GuestData `json:"data" validate:"required"`
	Metadata Metadata    `json:"metadata" validate:"required"`
}

type SingleActivityData struct {
	Data                             GuestData `json:"data" validate:"required"`
	CompetentAuthorityId_validated   string    `json:"competentAuthorityId_validated"`
	CompetentAuthorityName_validated string    `json:"competentAuthorityName_validated"`
	Metadata                         Metadata  `json:"metadata" validate:"required"`
}

type ActivityDataResponse struct {
	CA           string               `json:"ca"`
	ActivityData []SingleActivityData `json:"activityData"`
}

var logger = internal.CallNewLogger()

const activityDataTopic = "activity-data"

// ActivityDataPush godoc
//
//	@Summary		Submit activity data to the SDEP
//	@Schemes		https
//	@Description	Submit activity data to the SDEP
//	@Tags			str
//	@Accept			json
//	@Produce		json
//	@Param			ActivityData	body		ActivityData		true	"json formatted ActivityData"	example({"data":[{"numberOfGuests":3,"countryOfGuests":["ITA","NLD"],"temporal":{"startDateTime":"2024-07-21T17:32:28Z","endDateTime":"2024-07-25T17:32:28Z"},"address":{"street":"123 Main St","city":"Brussels","postalCode":"1000","country":"BEL"},"hostId":"placeholder-host-id","registrationNumber":"placeholder-registrationNumber"}],"metadata":{"platform":"booking.com","submissionDate":"2024-07-21T17:32:28Z","additionalProp1":{}}})
//	@Success		200				{array}		ActivityData		"delivered"
//	@Failure		400				{object}	BadRequestError		"Bad Request - Missing mandatory parameters"
//	@Failure		401				{object}	UnauthorizedError	"Unauthorized"	//	(Handled by OAuth2 middleware)
//	@Failure		422				{object}	BadRequestError		"Invalid data field values"
//	@Failure		503				{object}	InternalServerError	"Service unavailable"	
//	@Router			/str/activity-data [post]
//	@Security		OAuth2AccessCode[read]
func ActivityDataPush(ctx *gin.Context) {
	var newActivityData ActivityData

	validate := validator.New(validator.WithRequiredStructEnabled())

	// Call BindJSON to bind the received JSON
	if err := ctx.BindJSON(&newActivityData); err != nil {
        apiError := NewAPIError(
            "Missing mandatory input parameters",
            http.StatusBadRequest,
            "Invalid JSON format: " + err.Error(),
            ErrIDInvalidJSON,
        )
        ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError)
        return
    }

	// Validate the newActivityData format
	validationMessage, err := ValidateStruct(validate, newActivityData)
    if err != nil {
        logger.Error("Validation error", slog.String("error", err.Error()))
        apiError := NewAPIError(
            "Invalid data field values",
            http.StatusUnprocessableEntity,
            err.Error(),
            ErrIDValidationFailed,
        )
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, apiError)
        return
    }
	if validationMessage != "" {
		logger.Error("Validation error", slog.String("Message", validationMessage))
        apiError := NewAPIError(
            "Invalid data field values",
            http.StatusUnprocessableEntity,
            err.Error(),
            ErrIDValidationFailed,
        )
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, apiError)
        return
	}

	// Get the customClaims
	customClaims := middleware.GetCustomClaims(ctx)
	identity := Identity{OAuth2AppName: customClaims.Name, CA: strings.Split(customClaims.CA, ",")}

	identityJSON, _ := json.Marshal(identity)

	// Setup Kafka writer
	w := internal.Writer(activityDataTopic)

	// Use a derived context with timeout for pushing message
	kafkaCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	var messages []kafka.Message

	for _, data := range newActivityData.Data {
		singleActivityData := SingleActivityData{
			Data:     data,
			Metadata: newActivityData.Metadata,
		}

		// Call external API for RegistrationNumber and handle competent authority data
		if singleActivityData.Data.RegistrationNumber != "" {
			competentAuthority, err := callExternalAPI(singleActivityData.Data.RegistrationNumber)
			if err != nil {
				competentAuthority = CompetentAuthority{Name: "", Code: ""}
			}

			singleActivityData.CompetentAuthorityId_validated = competentAuthority.Code
			singleActivityData.CompetentAuthorityName_validated = competentAuthority.Name
		}

		// Validate the single activity data
		validationMessage, err := ValidateStruct(validate, singleActivityData)
		if err != nil {
			logger.Error("Validation error for single activity data", slog.String("error", err.Error()))
			apiError := NewAPIError(
				"Invalid data field values in activity data",
				http.StatusUnprocessableEntity,
				err.Error(),
				ErrIDValidationFailed,
			)
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, apiError)
			return
		}
		if validationMessage != "" {
			logger.Error("Validation message for single activity data", slog.String("message", validationMessage))
			apiError := APIError{
				Status: "error",
				Results: []ErrorEntry{
					{
						Message:      "Invalid data field values in activity data",
						Code:         http.StatusUnprocessableEntity, // 422
						ErrorUserMsg: validationMessage,
						ID:           "VALIDATION_004",
					},
				},
			}
			ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, apiError)
			return
		}

		singleActivityDataJSON, _ := json.Marshal(singleActivityData)

		messages = append(messages, kafka.Message{
			Headers: []kafka.Header{{Key: "Identity", Value: identityJSON}},
			Value:   singleActivityDataJSON,
		})
	}

	// Write messages to Kafka and handle errors
	err = w.WriteMessages(kafkaCtx, messages...)
	err = w.WriteMessages(kafkaCtx, messages...)
    if err != nil {
        logger.Error("Failed to write messages to Kafka", slog.String("error", err.Error()))
        apiError := NewAPIError(
            "Service temporarily unavailable",
            http.StatusServiceUnavailable,
            "Failed to process activity data. Please try again later.",
            ErrIDKafkaWriteFailed,
        )
        ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, apiError)
        return
    }

	if err := w.Close(); err != nil {
		logger.Error("Failed to close writer", slog.String("error", err.Error()))
		os.Exit(1)
	}

	ctx.JSON(http.StatusCreated, Status{Status: "delivered"})
}

// ActivityDataPull godoc
//
//	@Summary		Retrieve activity data submitted to the SDEP
//	@Schemes		https
//	@Description	Retrieve activity data submitted to the SDEP
//	@Param			limit	query	int	false	"Maximum number of records to return"
//	@Tags			ca
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		SingleActivityData
//	@Failure		400	{object}	BadRequestError		"Bad Request - Invalid limit parameter"
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"		//	(Handled by OAuth2 middleware)
//	@Failure		429	{object}	BadRequestError		"Too many requests"	//	(Could be enforced by rate-limiting middleware)
//	@Failure		503	{object}	InternalServerError	"Service unavailable"
//	@Router			/ca/activity-data [get]
//	@Security		OAuth2AccessCode[read]
func ActivityDataPull(ctx *gin.Context) {
	var activityData []SingleActivityData

	customClaims := middleware.GetCustomClaims(ctx)
	identity := Identity{OAuth2AppName: customClaims.Name, CA: strings.Split(customClaims.CA, ",")}

	r := internal.Reader(activityDataTopic, identity.OAuth2AppName)
	defer func() {
		if err := r.Close(); err != nil {
			logger.Error("Failed to close reader", slog.String("error", err.Error()))
		}
	}()

	kafkaCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	limit := ctx.DefaultQuery("limit", defaultLimit)
	maxMessages, err := strconv.Atoi(limit)
	if err != nil {
        apiError := NewAPIError(
            "Invalid limit parameter",
            http.StatusBadRequest,
            "The provided limit parameter is not a valid number",
            ErrIDInvalidParam,
        )
        ctx.JSON(http.StatusBadRequest, apiError)
        return
    }

	// Try to read messages with retries
	for i := 0; i < maxMessages; i++ {
		m, err := r.ReadMessage(kafkaCtx)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				logger.Info("No more messages available after reading", slog.Int("message", i))
				break
			}
			logger.Error("Failed to read message", slog.String("error", err.Error()))
			continue
		}

		var singleActivityData SingleActivityData
		if err := json.Unmarshal(m.Value, &singleActivityData); err != nil {
			logger.Error("Failed to unmarshal message", slog.String("error", err.Error()))
			continue
		}

		logger.Info("Processing message with CA validated",
			slog.String("message", singleActivityData.Data.CompetentAuthorityId_area),
			slog.String("message", singleActivityData.Data.CompetentAuthorityName_area),
			slog.String("message", singleActivityData.CompetentAuthorityId_validated))

		// Include message if no CA filter or empty CA
		if len(identity.CA) == 0 || identity.CA[0] == "" {
			activityData = append(activityData, singleActivityData)
			continue
		}

		for _, customCA := range identity.CA {
			competentAuthorityCode_validated := singleActivityData.CompetentAuthorityId_validated
			competentAuthorityCode_area := singleActivityData.Data.CompetentAuthorityId_area
			if competentAuthorityCode_area == customCA {
				activityData = append(activityData, singleActivityData)
			} else if competentAuthorityCode_validated == customCA {
				activityData = append(activityData, singleActivityData)
			}
		}
	}

	if len(activityData) == 0 {
		ctx.JSON(http.StatusOK, []SingleActivityData{})
		return
	}

	ctx.JSON(http.StatusOK, activityData)
}

func callExternalAPI(RegistrationNumber string) (CompetentAuthority, error) {
	apiURL := viper.GetString("CA_API_URL")
	apiKey := viper.GetString("CA_API_KEY")

	url := apiURL + RegistrationNumber
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CompetentAuthority{}, err
	}

	// Add API key to the Authorization header
	req.Header.Add("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return CompetentAuthority{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		logger.Error("Failed to call external API",
			slog.String("status", resp.Status),
			slog.Int("statusCode", resp.StatusCode))
		return CompetentAuthority{}, nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error("Failed to read response body", slog.String("error", err.Error()))
		return CompetentAuthority{}, err
	}

	var competentAuthority CompetentAuthority
	err = json.Unmarshal(body, &competentAuthority)
	if err != nil {
		logger.Error("Failed to unmarshal response body", slog.String("error", err.Error()))
		return CompetentAuthority{}, err
	}
	return competentAuthority, nil
}
