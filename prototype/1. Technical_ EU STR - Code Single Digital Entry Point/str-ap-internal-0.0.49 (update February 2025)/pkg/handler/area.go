package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/spf13/viper"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ShapefileType string

const (
	STRType  ShapefileType = "str"
	DataType ShapefileType = "data"
)

// ShapefileMetadata to store shapefile metadata
type ShapefileMetadata struct {
	ID                          string        `json:"id" example:"01J306Z5SJS720QPA5JHJQE7GD"`
	Name                        string        `json:"name" example:"Amsterdam.zip"`
	Type                        ShapefileType `json:"type" example:"str"`
	CompetentAuthorityId_area   string        `json:"competentAuthorityId_area" example:"competentAuthorityId_area" validate:"required"`
	CompetentAuthorityName_area string        `json:"competentAuthorityName_area" example:"competentAuthorityName_area" validate:"required"`
	Timestamp                   time.Time     `json:"timestamp" example:"2024-07-17T11:28:56.781531428Z"`
	// Add any other metadata fields you have in the JSON file
}

type ShapefileListResponse struct {
	Shapefiles []ShapefileMetadata `json:"shapefiles"`
}

type UploadResponse struct {
	Status                      string    `json:"status"`
	ID                          string    `json:"id"`
	Name                        string    `json:"name"`
	Timestamp                   time.Time `json:"timestamp"`
	CompetentAuthorityId_area   string    `json:"competentAuthorityId_area" example:"competentAuthorityId_area"`
	CompetentAuthorityName_area string    `json:"competentAuthorityName_area" example:"competentAuthorityName_area"`
}

func GetShapefiles(ctx *gin.Context, fileType ShapefileType) {
	var destDir string
	if fileType == STRType {
		destDir = viper.GetString("SHAPE_DEST") + "/STR"
	} else {
		destDir = viper.GetString("SHAPE_DEST") + "/Data"
	}
	if destDir == "" {
		apiError := NewAPIError(
            "Configuration error",
            http.StatusInternalServerError,
            "SHAPE_DEST environment variable not set",
            ErrIDServiceFailed,
        )
        ctx.JSON(http.StatusInternalServerError, apiError)
        return
	}

	// if destDir == "" {
	// 	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "SHAPE_DEST environment variable not set"})
	// 	return
	// }

	// Get a list of all files in the destination directory
	files, err := os.ReadDir(destDir)
	if err != nil {
		apiError := NewAPIError(
            "Internal server error",
            http.StatusInternalServerError,
            "Failed to read destination directory",
            ErrIDServiceFailed,
        )
        ctx.JSON(http.StatusInternalServerError, apiError)
        return
	}

	// Create a slice to store the shapefile metadata
	var shapefileList []ShapefileMetadata

	// Iterate over the files and add shapefile metadata to the list
	for _, file := range files {
		// Skip directories
		if file.IsDir() {
			continue
		}

		// Check if the file is a JSON file
		if filepath.Ext(file.Name()) == ".json" {
			// Build the file path
			filePath := filepath.Join(destDir, file.Name())

			// Read the JSON data from the file
			jsonData, err := os.ReadFile(filePath)
			if err != nil {
				log.Println("Failed to read JSON data from file:", err)
				continue
			}

			// Unmarshal the JSON data into a ShapefileMetadata struct
			var shapefileMetadata ShapefileMetadata
			err = json.Unmarshal(jsonData, &shapefileMetadata)
			if err != nil {
				log.Println("Failed to unmarshal JSON data:", err)
				continue
			}

			// Add the shapefile metadata to the list
			shapefileList = append(shapefileList, shapefileMetadata)
		}
	}

	// Return the shapefile list with metadata as the response
	ctx.JSON(http.StatusOK, shapefileList)
}

// DataAreaList_str godoc
//
//	@Summary		Get a list of available shapefiles (shapefiles where the STR regulation is applicable)
//	@Schemes		https
//	@Description	Retrieve a list of available STR shapefiles.
//	@Tags			str
//	@Produce		json
//	@Success		200	{array}		ShapefileMetadata
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"
//	@Router			/str/str-area [get]
//	@Security		OAuth2AccessCode[read]
func DataAreaList_str(ctx *gin.Context) {
	GetShapefiles(ctx, STRType)
}

// GetShapefiles-str godoc
//
//	@Summary		Get a list of available data-shapefiles (shapefiles for which you want to receive activity data for)
//	@Schemes		https
//	@Description	Retrieve a list of available Data shapefiles.
//	@Tags			str
//	@Produce		json
//	@Success		200	{array}		ShapefileMetadata
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"
//	@Router			/str/data-area [get]
//	@Security		OAuth2AccessCode[read]
func DataAreaList_data(ctx *gin.Context) {
	GetShapefiles(ctx, DataType)
}

func DownloadShapefile(ctx *gin.Context, fileType ShapefileType) {
	// Get the shapefile ID from the request URL parameter
	id := ctx.Param("id")
	// Load environment variables from .env file

	// Get the SHAPE_DEST environment variable
	var destDir string
	if fileType == STRType {
		destDir = viper.GetString("SHAPE_DEST") + "/STR"
	} else {
		destDir = viper.GetString("SHAPE_DEST") + "/Data"
	}

	if destDir == "" {
		apiError := NewAPIError(
            "Configuration error",
            http.StatusInternalServerError,
            "SHAPE_DEST environment variable not set",
            ErrIDServiceFailed,
        )
        ctx.JSON(http.StatusInternalServerError, apiError)
        return
	}

	// Build the shapefile path
	shapefilePath := filepath.Join(destDir, id)

	// Check if the shapefile exists
	_, err := os.Stat(shapefilePath)
	if os.IsNotExist(err) {
		apiError := NewAPIError(
            "Resource not found",
            http.StatusNotFound,
            "Shapefile not found",
            ErrIDResourceNotFound,
        )
        ctx.JSON(http.StatusNotFound, apiError)
        return
	}

	// Set the appropriate headers for the response
	ctx.Header("Content-Disposition", "attachment; filename="+id)
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.File(shapefilePath)
}

// DownloadShapefile godoc
//
//	@Summary		Download the shapefile for which you want to receive activity data. 
//	@Schemes		https
//	@Description	Retrieve the data shapefile by its ID and download it.
//	@Tags			str
//	@Produce		octet-stream
//	@Param			id	path		string	true	"Shapefile ID"
//	@Success		200	{file}		file
//	@Failure		400	{object}	BadRequestError		"Bad Request"
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"
//	@Failure		404	{object}	NotFoundError		"Shapefile not Found"
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"
//	@Router			/str/data-area/{luid} [get]
//	@Security		OAuth2AccessCode[read]
func DataAreaDownload(ctx *gin.Context) {
	DownloadShapefile(ctx, DataType)
}

// DownloadShapefile godoc
//
//	@Summary		Download the shapefile where the STR regulation is applicable. 
//	@Schemes		https
//	@Description	Retrieve the STR shapefile by its ID and download it.
//	@Tags			str
//	@Produce		octet-stream
//	@Param			id	path		string	true	"Shapefile ID"
//	@Success		200	{file}		file
//	@Failure		400	{object}	BadRequestError		"Bad Request"
//	@Failure		401	{object}	UnauthorizedError	"Unauthorized"
//	@Failure		404	{object}	NotFoundError		"Shapefile not Found"
//	@Failure		500	{object}	InternalServerError	"Internal Server Error"
//	@Router			/str/str-area/{luid} [get]
//	@Security		OAuth2AccessCode[read]
func StrAreaDownload(ctx *gin.Context) {
	DownloadShapefile(ctx, STRType)
}

func handleAreaUpload(ctx *gin.Context, fileType ShapefileType) {
    var destDir string
    if fileType == STRType {
        destDir = viper.GetString("SHAPE_DEST") + "/STR"
    } else {
        destDir = viper.GetString("SHAPE_DEST") + "/Data"
    }

    if destDir == "" {
        apiError := NewAPIError(
            "Configuration error",
            http.StatusInternalServerError,
            "Destination directory environment variable not set",
            ErrIDServiceFailed,
        )
        ctx.JSON(http.StatusInternalServerError, apiError)
        return
    }

    // Get the uploaded file from the request
    file, err := ctx.FormFile("file")
    if err != nil {
        apiError := NewAPIError(
            "Invalid request",
            http.StatusBadRequest,
            "Failed to get file from request",
            ErrIDInvalidJSON,
        )
        ctx.JSON(http.StatusBadRequest, apiError)
        return
    }

    // Get the competent authority fields from the request
    competentauthorityidArea := ctx.PostForm("competentAuthorityId")
    competentauthoritynameArea := ctx.PostForm("competentAuthorityName")

    if competentauthorityidArea == "" || competentauthoritynameArea == "" {
        apiError := NewAPIError(
            "Missing required fields",
            http.StatusBadRequest,
            "Competent Authority Name and Code are required",
            ErrIDMissingField,
        )
        ctx.JSON(http.StatusBadRequest, apiError)
        return
    }

    // Open the uploaded file
    src, err := file.Open()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
        return
    }
    defer func(src multipart.File) {
        err := src.Close()
        if err != nil {
            log.Printf("Error closing source file: %v", err)
        }
    }(src)

    id := ulid.Make().String()
    destPath := filepath.Join(destDir, id)
    // Check if the destination directory exists
    if err := os.MkdirAll(destDir, 0755); err != nil {
        os.Create(destPath)
    }

    // Create the destination file
    dst, err := os.Create(destPath)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create destination file"})
        return
    }
    defer func(dst *os.File) {
        err := dst.Close()
        if err != nil {
            log.Printf("Error closing destination file: %v", err)
        }
    }(dst)

    // Copy the contents of the uploaded file to the destination file
    _, err = io.Copy(dst, src)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Store the shapefile info
    shapefileInfo := ShapefileMetadata{
        ID:                          id,
        Name:                        file.Filename,
        Type:                        fileType,
        Timestamp:                   time.Now(),
        CompetentAuthorityId_area:   competentauthorityidArea,
        CompetentAuthorityName_area: competentauthoritynameArea,
    }

    // Create the JSON file
    jsonFilePath := filepath.Join(destDir, id+".json")
    jsonFile, err := os.Create(jsonFilePath)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create JSON file"})
        return
    }
    defer func(jsonFile *os.File) {
        err := jsonFile.Close()
        if err != nil {
            log.Printf("Error closing JSON file: %v", err)
        }
    }(jsonFile)

    // Marshal the shapefile info to JSON
    jsonData, err := json.Marshal(shapefileInfo)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON data"})
        return
    }

    // Write the JSON data to the file
    _, err = jsonFile.Write(jsonData)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write JSON data to file"})
        return
    }

    // Build the response JSON
    response := gin.H{
        "status":                      "Successfully uploaded",
        "ID":                          id,
        "Name":                        file.Filename,
        "Type":                        fileType,
        "CompetentAuthorityId_area":   shapefileInfo.CompetentAuthorityId_area,
        "CompetentAuthorityName_area": shapefileInfo.CompetentAuthorityName_area,
        "Timestamp":                   shapefileInfo.Timestamp,
    }

    ctx.JSON(http.StatusOK, response)
}

// StrAreaUpload godoc
//
//	@Summary		Upload a new STR area shapefile
//	@Schemes		https
//	@Description	Upload a new STR area shapefile to the server.
//	@Tags			ca
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file					formData	file	true	"Shapefile to upload"
//	@Param			competentAuthorityId	formData	string	true	"Id of the Competent Authority"
//	@Param			competentAuthorityName	formData	string	true	"Name of the Competent Authority"
//	@Success		200						{object}	UploadResponse
//	@Failure		400						{object}	BadRequestError		"Bad Request - Missing mandatory parameters"
//	@Failure		401						{object}	UnauthorizedError	"Unauthorized"
//	@Failure		422						{object}	BadRequestError		"Invalid file format or data"
//	@Failure		503						{object}	InternalServerError	"Service unavailable"
//	@Router			/ca/str-area [post]
//	@Security		OAuth2AccessCode[write]
func StrAreaUpload(ctx *gin.Context) {
	handleAreaUpload(ctx, STRType)
}

// DataAreaUpload godoc
//
//	@Summary		Upload a new data area shapefile
//	@Schemes		https
//	@Description	Upload a new data area shapefile to the server.
//	@Tags			ca
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file					formData	file	true	"Shapefile to upload"
//	@Param			competentAuthorityId	formData	string	true	"Id of the Competent Authority"
//	@Param			competentAuthorityName	formData	string	true	"Name of the Competent Authority"
//	@Success		200						{object}	UploadResponse
//	@Failure		400						{object}	BadRequestError		"Bad Request - Missing mandatory parameters"
//	@Failure		401						{object}	UnauthorizedError	"Unauthorized"
//	@Failure		422						{object}	BadRequestError		"Invalid file format or data"
//	@Failure		503						{object}	InternalServerError	"Service unavailable"
//	@Router			/ca/data-area [post]
//	@Security		OAuth2AccessCode[write]
func DataAreaUpload(ctx *gin.Context) {
	handleAreaUpload(ctx, DataType)
}

func AreaDelete(ctx *gin.Context, fileType ShapefileType) {
	var destDir string
	if fileType == STRType {
		destDir = viper.GetString("SHAPE_DEST") + "/STR"
	} else {
		destDir = viper.GetString("SHAPE_DEST") + "/Data"
	}
	if destDir == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "SHAPE_DEST environment variable not set"})
		return
	}

	// Get the LUID from the URL
	luid := ctx.Param("luid")
	if luid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid LUID format"})
		return
	}

	// Build the shapefile and JSON file paths
	shapefilePath := filepath.Join(destDir, luid)
	jsonFilePath := filepath.Join(destDir, luid+".json")

	// Check if the shapefile exists
	if _, err := os.Stat(shapefilePath); os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Shapefile not found"})
		return
	}

	// Delete the shapefile
	if err := os.Remove(shapefilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete shapefile"})
		return
	}

	// Delete the JSON file
	if err := os.Remove(jsonFilePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete JSON file"})
		return
	}

	// Build the response JSON
	response := gin.H{
		"status": "Successfully deleted",
	}

	ctx.JSON(http.StatusOK, response)
}

// AreaDelete godoc
//
//	@Summary		Delete a STR area shapefile
//	@Schemes		https
//	@Description	Delete the STR shapefile from the server based on the LUID.
//	@Tags			ca
//	@Produce		json
//	@Param			luid	path		string	true	"LUID of the shapefile to delete"
//	@Success		200		{object}	DeleteResponse
//	@Failure		400		{object}	BadRequestError		"Bad Request"
//	@Failure		401		{object}	UnauthorizedError	"Unauthorized"
//	@Failure		404		{object}	NotFoundError		"Not Found"
//	@Failure		500		{object}	InternalServerError	"Internal Server Error"
//	@Router			/ca/str-area/{luid} [delete]
//	@Security		OAuth2AccessCode[write]
func DataAreaDelete_STR(ctx *gin.Context) {
	AreaDelete(ctx, STRType)
}

// AreaDelete godoc
//
//	@Summary		Delete a data area shapefile
//	@Schemes		https
//	@Description	Delete a data shapefile from the server based on the LUID.
//	@Tags			ca
//	@Produce		json
//	@Param			luid	path		string	true	"LUID of the shapefile to delete"
//	@Success		200		{object}	DeleteResponse
//	@Failure		400		{object}	BadRequestError		"Bad Request"
//	@Failure		401		{object}	UnauthorizedError	"Unauthorized"
//	@Failure		404		{object}	NotFoundError		"Not Found"
//	@Failure		500		{object}	InternalServerError	"Internal Server Error"
//	@Router			/ca/data-area/{luid} [delete]
//	@Security		OAuth2AccessCode[write]
func DataAreaDelete_Data(ctx *gin.Context) {
	AreaDelete(ctx, DataType)
}