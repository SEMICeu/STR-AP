package handler

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	kafka "github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"io"
	"net/http"
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
	Address           Address             `json:"address"`  
	CompetentAuthorityId_area           string              `json:"competentAuthorityId_area" example:"competentAuthorityId_area"`
	CompetentAuthorityName_area            string              `json:"competentAuthorityName_area" example:"competentAuthorityName_area"`
	CountryOfGuests   []string            `json:"countryOfGuests" example:"ITA,NLD"`  
	HostID            string              `json:"hostId" example:"placeholder-host-id"`  
	NumberOfGuests    uint                `json:"numberOfGuests" example:"3"` 
	Temporal          Temporal            `json:"temporal"`   
	URL            string              `json:"URL" example:"placeholder-URL"`   				
	RegistrationNumber            string              `json:"registrationNumber" example:"placeholder-registrationNumber" validate:"required"`
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
	Data     GuestData `json:"data" validate:"required"`
	CompetentAuthorityId_validated string `json:"competentAuthorityId_validated"`  
	CompetentAuthorityName_validated string `json:"competentAuthorityName_validated"`   
	Metadata Metadata  `json:"metadata" validate:"required"` 
}  
  
type ActivityDataResponse struct {  
	CA           string              `json:"ca"`  
	ActivityData []SingleActivityData `json:"activityData"`  
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
//	@Param			ActivityData	body		ActivityData		true	"json formatted ActivityData"	example({"data":[{"numberOfGuests":3,"countryOfGuests":["ITA","NLD"],"temporal":{"startDateTime":"2024-07-21T17:32:28Z","endDateTime":"2024-07-25T17:32:28Z"},"address":{"street":"123 Main St","city":"Brussels","postalCode":"1000","country":"BEL"},"hostId":"placeholder-host-id","registrationNumber":"placeholder-registrationNumber"}],"metadata":{"platform":"booking.com","submissionDate":"2024-07-21T17:32:28Z","additionalProp1":{}}})	
//	@Success		200				{array}		ActivityData		"delivered"
//	@Failure		400				{object}	BadRequestError		"Bad Request"	
//	@Failure		401				{object}	UnauthorizedError	"Unauthorized"	
//	@Router			/str/activity-data [post]  
//	@Security		OAuth2AccessCode[read]  
func ActivityDataPush(ctx *gin.Context) {  
	var newActivityData ActivityData  
  
	validate := validator.New(validator.WithRequiredStructEnabled())  
  
	statusWrongDataFormat := Status{Status: "Wrong data format!"}  
  
	// Call BindJSON to bind the received JSON  
	if err := ctx.BindJSON(&newActivityData); err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, statusWrongDataFormat)  
		return  
	}  
  
	// Validate the newActivityData format  
	validationMessage, err := ValidateStruct(validate, newActivityData) // Renamed variable for clarity  
	if err != nil {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: err.Error()})  
		return  
	}  
	if validationMessage != "" {  
		ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: validationMessage})  
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
  
	messages := []kafka.Message{}  

	for _, data := range newActivityData.Data {  
		singleActivityData := SingleActivityData{  
			Data:     data,  
			Metadata: newActivityData.Metadata,  
		}  
  
        // Validate the data  
        if validationError, err := ValidateStruct(validate, singleActivityData); err != nil || validationError != "" {    
            ctx.AbortWithStatusJSON(http.StatusBadRequest, HTTPError{Code: http.StatusBadRequest, Message: err.Error() + validationError})    
            return    
        }    
 
        // Call external API for RegistrationNumber and handle competent authority data  
        competentAuthority, err := callExternalAPI(singleActivityData.Data.RegistrationNumber)    
       
		if err != nil {    
            competentAuthority = CompetentAuthority{Name: "", Code: ""}    
        }    
		
        singleActivityData.CompetentAuthorityId_validated= competentAuthority.Code
        singleActivityData.CompetentAuthorityName_validated = competentAuthority.Name
 
        singleActivityDataJSON, _ := json.Marshal(singleActivityData)

        messages = append(messages, kafka.Message{    
            Headers: []kafka.Header{{Key: "Identity", Value: identityJSON}},    
            Value:   singleActivityDataJSON,    
        })
    }	
 
    // Write messages to Kafka and handle errors  
    if err := w.WriteMessages(kafkaCtx, messages...); err != nil {    
        ctx.AbortWithStatusJSON(http.StatusInternalServerError, statusWrongDataFormat)    
        internal.Fatalf("Failed to write messages: %s", err)    
        return    
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
//	@Param			limit	query	int	false	"Maximum number of records to return"	
//	@Tags			id 
//	@Accept			json  
//	@Produce		json  
//	@Success		200	{array}		SingleActivityData	
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"			
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"	
//	@Router			/ca/activity-data [get]  
//	@Security		OAuth2AccessCode[read]  
func ActivityDataPull(ctx *gin.Context) {
    var activityData []SingleActivityData


    // Get the customClaims
    customClaims := middleware.GetCustomClaims(ctx)
    identity := Identity{OAuth2AppName: customClaims.Name, CA: strings.Split(customClaims.CA, ",")}


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

        if len(identity.CA) == 0 || identity.CA[0] == "" {    
            activityData = append(activityData, singleActivityData)
        } else {    
            for _, customCA := range identity.CA {      
                // apply filter on activity data      
                competentAuthorityCode_validated := singleActivityData.CompetentAuthorityId_validated
				competentAuthorityCode_area := singleActivityData.Data.CompetentAuthorityId_area
				if competentAuthorityCode_area == customCA {    
                    activityData = append(activityData, singleActivityData)
				} else if competentAuthorityCode_validated == customCA {    
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
	defer resp.Body.Close()  
  
	if resp.StatusCode != http.StatusOK { 
		internal.Errorf("Failed to get competent authority data, API http return code: %v", resp.StatusCode) 
		return CompetentAuthority{}, nil
	}  
  
	body, err := io.ReadAll(resp.Body)  
	if err != nil { 
		internal.Errorf("Unable to read response body") 
		return CompetentAuthority{}, err  
	}  
  
	var competentAuthority CompetentAuthority  
	err = json.Unmarshal(body, &competentAuthority) 
	if err != nil {  
		return CompetentAuthority{}, err  
	}  
	return competentAuthority, nil  
}