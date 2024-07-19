package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	kafka "github.com/segmentio/kafka-go"
	"net/http"
	"str/internal"
	"str/pkg/middleware"
	"strconv"
	"strings"
	"time"
)

// Nested structures

type Temporal struct {
	StartDateTime time.Time `json:"startDateTime" example:"2024-07-21T17:32:28Z"`
	EndDateTime   time.Time `json:"endDateTime" example:"2024-07-25T17:32:28Z"`
}

type GuestData struct {
	NumberOfGuests  uint     `json:"numberOfGuests" example:"3"`
	CountryOfGuests []string `json:"countryOfGuests" example:"ITA,NLD"`
	Temporal        Temporal `json:"temporal"`
	Address         Address  `json:"address"`
	HostID          string   `json:"hostId" example:"placeholder-host-id"`
	UnitID          string   `json:"unitId" example:"placeholder-unit-id" validate:"required"`
	AreaID          string   `json:"areaId" example:"placeholder-area-id"`
}

type ActivityData struct {
	Data     []GuestData `json:"data" validate:"required"`
	Metadata Metadata    `json:"metadata" validate:"required"`
}

type SingleActivityData struct {
	Data     GuestData `json:"data" validate:"required"`
	Metadata Metadata  `json:"metadata" validate:"required"`
}

const activityDataTopic = "activity-data"

// ActivityDataPush godoc
//
//	@Summary		Submit activity data to the SDEP
//	@Schemes		https
//	@Description	Submit activity data to the SDEP
//	@Tags			str
//	@Accept			json
//	@Produce		json
//	@Param			ActivityData	body		ActivityData	true	"json formatted ActivityData"	example({"data":[{"numberOfGuests":3,"countryOfGuests":["ITA","NLD"],"temporal":{"startDateTime":"2024-07-21T17:32:28Z","endDateTime":"2024-07-25T17:32:28Z"},"address":{"street":"123 Main St","city":"Brussels","postalCode":"1000","country":"BEL"},"hostId":"placeholder-host-id","unitId":"placeholder-unit-id","areaId":"placeholder-area-id"}],"metadata":{"platform":"booking.com","submissionDate":"2024-07-21T17:32:28Z","additionalProp1":{}}})
//	@Success		201				{object}	Status
//	@Failure		400				{object}	HTTPError
//	@Router			/str/activity-data [post]
//	@Security		OAuth2AccessCode[read]
func ActivityDataPush(ctx *gin.Context) {
	var newActivityData ActivityData

	validate = validator.New(validator.WithRequiredStructEnabled())

	statusWrongDataFormat := Status{Status: "Wrong data format!"}

	// Call BindJSON to bind the received JSON
	if err := ctx.BindJSON(&newActivityData); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, statusWrongDataFormat)
		return
	}

	// Validate the newActivityData format
	validationError, err := ValidateStruct(validate, newActivityData)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			internal.Fatalf(err.Error())
		}
	}
	if validationError != "" {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: validationError})
		return
	}

	// Get the customClaims
	customClaims := middleware.GetCustomClaims(ctx)
	identity := Identity{OAuth2AppName: customClaims.Name, CA: customClaims.CA}

	identityJSON, _ := json.Marshal(identity)

	// Setup Kafka writer
	w := internal.Writer(activityDataTopic)

	// Use a derived context with timeout for pushing message
	kafkaCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	messages := []kafka.Message{}

	for _, data := range newActivityData.Data {
		singleActivityData := SingleActivityData{
			Data:     data,
			Metadata: newActivityData.Metadata,
		}

		// Validate the data
		validationError, err := ValidateStruct(validate, singleActivityData)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				internal.Fatalf(err.Error())
			}
		}
		if validationError != "" {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: validationError})
			return
		}

		singleActivityDataJSON, _ := json.Marshal(singleActivityData)

		messages = append(messages, kafka.Message{
			Headers: []kafka.Header{{Key: "Identity", Value: []byte(identityJSON)}},
			Value:   []byte(singleActivityDataJSON),
		})
	}

	err = w.WriteMessages(kafkaCtx, messages...)
	if err != nil {
		internal.Errorf("Failed to write messages: %s", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, statusWrongDataFormat)
	}

	if err := w.Close(); err != nil {
		internal.Fatalf("Failed to close writer: %s", err)
	}

	statusDelivered := Status{Status: "delivered"}
	ctx.JSON(http.StatusCreated, statusDelivered)
}

// ActivityDataPull godoc
//
//	@Summary		Retrieve activity data submitted to the SDEP
//	@Schemes		https
//	@Description	Retrieve activity data submitted to the SDEP
//	@Param			limit	query	int	false	"limit number of records returned"	minimum(1)	maximum(100)
//	@Tags			ca
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Status
//	@Failure		400	{object}	HTTPError
//	@Router			/ca/activity-data [get]
//	@Security		OAuth2AccessCode[read]
func ActivityDataPull(ctx *gin.Context) {
	var activityData []SingleActivityData

	// Get the customClaims
	customClaims := middleware.GetCustomClaims(ctx)
	identity := Identity{OAuth2AppName: customClaims.Name, CA: customClaims.CA}

	// todo handle request parameter, to wipe out the groupID, so we get all the data from the beginning
	// Setup Kafka reader
	r := internal.Reader(activityDataTopic, identity.OAuth2AppName)

	kafkaCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := ctx.DefaultQuery("limit", defaultLimit)
	// Define the maximum number of messages to read
	maxMessages, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	for i := 0; i < maxMessages; i++ {
		var singleActivityData SingleActivityData
		m, err := r.ReadMessage(kafkaCtx)
		if err != nil {
			break
		}

		err = json.Unmarshal(m.Value, &singleActivityData)

		if identity.CA == "" {
			activityData = append(activityData, singleActivityData)
		} else {
			// apply filter on activity data
			parts := strings.Split(singleActivityData.Data.AreaID, "-")
			if len(parts) == 3 {
				middleElement := parts[1]
				if middleElement == identity.CA {
					activityData = append(activityData, singleActivityData)
				}
			}
		}
	}

	ctx.JSON(http.StatusOK, activityData)

	if err := r.Close(); err != nil {
		internal.Fatalf("failed to close reader: %s", err)
	}
}
