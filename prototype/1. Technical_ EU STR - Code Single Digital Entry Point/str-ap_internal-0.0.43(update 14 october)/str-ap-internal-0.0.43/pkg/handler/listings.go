package handler  
  
import (  
	"context"  
	"encoding/json"  
	"github.com/gin-gonic/gin"  
	"github.com/go-playground/validator/v10"  
	"github.com/segmentio/kafka-go"  
	"net/http"  
	"str/internal"  
	"str/pkg/middleware"  
	"strconv"  
	"strings"  
	"time"  
)  
  
type Unit struct {  
	Description   string  `json:"description"`  
	FloorLevel    string  `json:"floorLevel"`  
	Address       Address `json:"address"`  
	ObtainedAuth  bool    `json:"obtainedAuth"`  
	SubjectToAuth bool    `json:"subjectToAuth"`  
	NumberOfRooms int     `json:"numberOfRooms"`  
	Occupancy     int     `json:"occupancy"`  
	Purpose       string  `json:"purpose"`  
	Type          string  `json:"type"`  
	URL           string  `json:"url"`  
}  
  
type RegistrationData struct {  
	RegistrationNumber string `json:"registrationNumber"`  
	Unit               Unit   `json:"Unit"` 
	CompetentAuthorityId_area           string              `json:"competentAuthorityId_area" example:"competentAuthorityId_area"`
	CompetentAuthorityName_area            string              `json:"competentAuthorityName_area" example:"competentAuthorityName_area"`
}  
  
type ListingsData struct {  
	Data     []RegistrationData `json:"data" validate:"required"`  
	Metadata Metadata           `json:"metadata" validate:"required"`  
}  
  
type SingleListingData struct {  
	Data     RegistrationData `json:"data" validate:"required"`  
	Metadata Metadata         `json:"metadata" validate:"required"`  
}  
  
const (  
	listingsTopic = "listings"  
	defaultLimit  = "100"  
)  
  
// ListingsPush godoc  
//  
//	@Summary		submit listing(s) to SDEP  
//	@Schemes		https  
//	@Description	submit listing(s) to SDEP  
//	@Tags			str  
//	@Accept			json  
//	@Produce		json  
//	@Param			ListingsData	body		ListingsData		true			"json formatted ListingsData"	example({"data":[{"registrationNumber":"1234","Unit":{"description":"string","floorLevel":"string","address":{"street":"Culliganlaan 5","city":"Diegem","postalCode":"1831","country":"BEL"},"obtainedAuth":true,"subjectToAuth":true,"numberOfRooms":0,"occupancy":0,"purpose":"string","type":"string","url":"STR-Platform.com/1234"}}],"metadata":{"platform":"STR-Platform"}}))	
//	@Success		201				{object}	Created				"Created"		
//	@Failure		400				{object}	BadRequestError		"Bad Request"	
//	@Failure		401				{object}	UnauthorizedError	"Unauthorized"	
//	@Router			/str/listings [post]  
//	@Security		OAuth2AccessCode[read]  
func ListingsPush(ctx *gin.Context) {  
	var newListings ListingsData  
  
	validate := validator.New(validator.WithRequiredStructEnabled()) // Corrected missing `:=`  
  
	statusWrongDataFormat := Status{Status: "Wrong data format!"}  
  
	// Call BindJSON to bind the received JSON  
	if err := ctx.BindJSON(&newListings); err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, statusWrongDataFormat)  
		return  
	}  
  
	// Validate the newListings format  
	validationError, err := ValidateStruct(validate, newListings)  
	if err != nil {  
		// Simplified the error check  
		internal.Fatalf(err.Error())  
	}  
	if validationError != "" {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: validationError})  
		return  
	}  
  
	// Get the customClaims  
	customClaims := middleware.GetCustomClaims(ctx)  
	identity := Identity{OAuth2AppName: customClaims.Name, CA: strings.Split(customClaims.CA, ",")}  
  
	identityJSON, _ := json.Marshal(identity)  
  
	// Setup Kafka writer  
	w := internal.Writer(listingsTopic)  
  
	// Use a derived context with timeout for pushing message  
	kafkaCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)  
	defer cancel()  
  
	messages := []kafka.Message{}  
  
	for _, data := range newListings.Data {  
		singleListingData := SingleListingData{  
			Data:     data,  
			Metadata: newListings.Metadata,  
		}  
  
		// Validate the data  
		validationError, err := ValidateStruct(validate, singleListingData)  
		if err != nil {  
			// Simplified the error check  
			internal.Fatalf(err.Error())  
		}  
		if validationError != "" {  
			ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: validationError})  
			return  
		}  
  
		singleListingDataJSON, _ := json.Marshal(singleListingData)  
  
		messages = append(messages, kafka.Message{  
			Headers: []kafka.Header{{Key: "Identity", Value: []byte(identityJSON)}},  
			Value:   []byte(singleListingDataJSON),  
		})  
	}  
  
	err = w.WriteMessages(kafkaCtx, messages...)  
	if err != nil {  
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, statusWrongDataFormat)  
		internal.Fatalf("Failed to write messages: %s", err)  
	}  
  
	if err := w.Close(); err != nil {  
		internal.Fatalf("Failed to close writer: %s", err)  
	}  
  
	statusDelivered := Status{Status: "delivered"}  
	ctx.JSON(http.StatusCreated, statusDelivered)  
}  
  
// ListingsPull godoc  
//  
//	@Summary		Retrieve listings submitted to the SDEP  
//	@Schemes		https  
//	@Description	Retrieve listings submitted to the SDEP  
//	@Param			limit	query	int	false	"limit number of records returned"	minimum(1)	maximum(100)	
//	@Tags			ca  
//	@Produce		json  
//	@Success		200	{object}	ListingsData		
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"			
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"	
//	@Router			/ca/listings [get]  
//	@Security		OAuth2AccessCode[read]  
func ListingsPull(ctx *gin.Context) {  
	var listingsData []SingleListingData  
  
	// Get the customClaims  
	customClaims := middleware.GetCustomClaims(ctx)  
	identity := Identity{OAuth2AppName: customClaims.Name, CA: strings.Split(customClaims.CA, ",")}  
  
	// Setup Kafka reader  
	r := internal.Reader(listingsTopic, identity.OAuth2AppName)  
  
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
		var singleListingData SingleListingData  
		m, err := r.ReadMessage(kafkaCtx)  
		if err != nil {  
			break  
		}  
  
		err = json.Unmarshal(m.Value, &singleListingData)  
		if err != nil {  
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal JSON"})  
			return  
		}  
		listingsData = append(listingsData, singleListingData)  
	}  
  
	ctx.JSON(http.StatusOK, listingsData)  
  
	if err := r.Close(); err != nil {  
		internal.Fatalf("failed to close reader: %s", err)  
	}  
}  
