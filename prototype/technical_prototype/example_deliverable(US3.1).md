| Disclaimer  |
|-----------|
| This Readme was prepared for DG Grow by PwC EU Services. The views expressed in this report are purely those of the authors and may not, in any circumstances, be interpreted as stating an official position of the European Commission. The European Commission does not guarantee the accuracy of the information included in this report, nor does it accept any responsibility for any use thereof. Reference herein to any specific products, specifications, process, or service by trade name, trademark, manufacturer, or otherwise, does not necessarily constitute or imply its endorsement, recommendation, or favouring by the European Commission. All care has been taken by the author to ensure that s/he has obtained, where necessary, permission to use any parts of manuscripts including illustrations, maps, and graphs, on which intellectual property rights already exist from the titular holder(s) of such rights or from her/his or their legal representative.|

# 1. Introduction  
This document outlines the architecture and implementation details for sharing activity data from STR platforms to the SDEP. It covers the ability to receive and process activity data, ensuring secure and efficient data transmission using Kafka. The detailed Swagger documentation for the POST /api/v0/str/activity-data endpoint is included to provide a comprehensive guide for developers.

# 2. Architecture Overview

## 2.1. Components
- API Gateway: Routes incoming requests to appropriate microservices.
- SDEP: Handles the processing and retrieval of activity data.
- Kafka: Message broker for reliable data transmission.


## 2.2 Technologies 

- Programming Language: Go
- API Documentation: Swagger/OpenAPI
- Messaging: Kafka
- Deployment: Kubernetes, HELM

## 2.3 FLOW
**1. Data Submission by STR Platforms:**
- The STR platform sends a POST request to the SDEP endpoint /api/v0/str/activity-data with the activity data.
  
**2. Decompose data object**
  
**3. Data Validation:**
  
Example Validation Checks:
- Ensure unitId is present and valid.
- Check if numberOfGuests is a positive integer.
- Validate the date format in stayDuration.
- Ensure countryOfResidence matches a recognized country code.
- Ensure aAreaID matches a known areaIid.
  
**4. Kafka Topic for Activity Data:**
A dedicated Kafka topic, say activity-data, is created to handle the activity data submissions.
Each record in this topic represents an entry of activity data submitted by an STR platform.

## 2.4 Endpoint Documentation 

**Endpoint:** POST /api/v0/str/activity-data

**Description:**  Allows an STR platform to submit activity data, which includes details like the number of guests, the period of stay, and the location of the unit.

**Request:**
- URL: /api/v0/str/activity-data
- Method: POST
- Content-Type: application/json
- Authentication: Bearer Token (JWT)
  
**Request Body:**

```json
{
  "data": [
    {
      "numberOfGuests": 3,
      "countryOfGuests": [
        "ITA",
        "NLD"
      ],
      "temporal": {
        "startDateTime": "2024-07-21T17:32:28Z",
        "endDateTime": "2024-07-25T17:32:28Z"
      },
      "address": {
        "street": "123 Main St",
        "city": "Brussels",
        "postalCode": "1000",
        "country": "BEL"
      },
      "hostId": "placeholder-host-id",
      "unitId": "placeholder-unit-id",
      "areaId": "placeholder-area-id"
    }
  ],
  "metadata": {
    "platform": "booking.com",
    "submissionDate": "2024-07-21T17:32:28Z",
    "additionalProp1": {}
  }
}
```

**Request Parameters:**
- hostId: The ID of the host providing the accommodation.
- unitId: The ID of the unit being rented.
- numberOfGuests: The number of guests staying at the unit.
- temporal: The period during which the unit is rented, with startDate and endDate.
- address: The physical address of the unit.
- areaId: The ID of the Area of the relevant CA 

**Responses:**

**201 OK**
```json
{
  "status": "success",
  "message": "Activity data submitted successfully."
}
```
**400 Bad Request**
```json
{
  "status": "error",
  "message": "Invalid input data."
}
```
**401 Unauthorized**
```json
{
  "status": "error",
  "message": "Unauthorized access."
}
```
**500 Internal Server Error**
```json
{
  "status": "error",
  "message": "An unexpected error occurred."
```



