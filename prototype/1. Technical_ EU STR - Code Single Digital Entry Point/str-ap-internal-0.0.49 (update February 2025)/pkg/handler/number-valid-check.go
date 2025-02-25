package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

// CallExternalAPI_check godoc
//
//	@Summary		Registration Number Validation check
//	@Schemes		https
//	@Description	Registration Number Validation check
//	@Tags			str
//	@Accept			json
//	@Produce		json
//	@Param			registrationNumber	body		string				true	"Registration Number"
//	@Success		200					{object}	Numbervalid			"Success"
//	@Failure		400					{object}	BadRequestError		"Bad Request - Missing registration number"
//	@Failure		401					{object}	UnauthorizedError	"Unauthorized"
//	@Failure		422					{object}	BadRequestError		"Invalid registration number format"
//	@Failure		503					{object}	InternalServerError	"Service unavailable"
//	@Router			/str/number-valid-check/{registrationNumber} [get]
//	@Security		OAuth2AccessCode[read]
func CallExternalAPI_check(ctx *gin.Context) {
	apiURL := viper.GetString("CA_API_URL_check")
	apiKey := viper.GetString("CA_API_KEY")

	regnnumber := ctx.Param("rnumb")
	if regnnumber == "" {
		apiError := NewAPIError(
            "Missing registration number",
            http.StatusBadRequest,
            "Registration number is required",
            ErrIDMissingField,
        )
        ctx.AbortWithStatusJSON(http.StatusBadRequest, apiError)
        return
	}

	body, err := fetchAPIResponse(apiURL, regnnumber, apiKey)
	if err != nil {
		apiError := NewAPIError(
            "Service temporarily unavailable",
            http.StatusServiceUnavailable,
            "Failed to validate registration number. Please try again later.",
            ErrIDExternalAPIFailed,
        )
        ctx.AbortWithStatusJSON(http.StatusServiceUnavailable, apiError)
        return
	}

	// Check if the response is empty or invalid
	if len(body) == 0 {
		apiError := NewAPIError(
            "Invalid registration number format",
            http.StatusUnprocessableEntity,
            "The provided registration number format is invalid",
            ErrIDFormatValidation,
        )
        ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, apiError)
        return
	}

	ctx.Data(http.StatusOK, "application/json", body)
}

// fetchAPIResponse sends a GET request to the API and returns the response body
func fetchAPIResponse(apiURL, regnnumber, apiKey string) ([]byte, error) {
	url := apiURL + regnnumber
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Api-Key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}