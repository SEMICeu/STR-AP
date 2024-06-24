# EU STR - Single Digital Entry Point
Base URL: eu-str.sdep-pilot.eu/api/v0

Swagger: https://eu-str.sdep-pilot.eu/swagger/index.html#/

A gateway for teh electronic transmission of data between online short-term rental platforms and competent authorities, ensuring priority of development is: 1. listings, 2. orders, 3. activity, 4. area. 

To obtain API credentials, please contact: wouter.travers@pwc.com via email. 

## Setup
This notebook is a collection of `curl` commands that can be executed in the terminal to validate the API endpoints\.
Besides `curl` also `jq` is used for formatting purposes\. 
### client credentials
The `bw` command is the Bitwarden password manager CLI command\.
Replace the CLIENT\_ID value\, with your personal CLIENT\_ID\. And provison your CLIENT\_SECRET
```warp-runnable-command
CLIENT_ID=plQ8P4mNxEoKFLY0NnYvfkO1Ak9HozjG \
CLIENT_SECRET=$(bw get password plQ8P4mNxEoKFLY0NnYvfkO1Ak9HozjG)
```
The client SECRET can also be captured via CLI read\:
```warp-runnable-command
read CLIENT_SECRET
```
Get an OAUTH token \( from the `/token` endpoint\)
```warp-runnable-command
# Compose the JSON string using jq
DATA=$(jq -n \
          --arg client_id "$CLIENT_ID" \
          --arg client_secret "$CLIENT_SECRET" \
          --arg audience "https://str.eu" \
          --arg grant_type "client_credentials" \
          '{client_id: $client_id, client_secret: $client_secret, audience: $audience, grant_type: $grant_type}')

# Get the token          
TOKEN=$(curl -s --request POST \
  --url https://tt-dp-dev.eu.auth0.com/oauth/token \
  --header 'content-type: application/json' \
  --data $DATA | jq -r .access_token )
```
### API host
```warp-runnable-command
HOST=eu-str.sdep-pilot.eu

```
## Test API
### Unauthenticated health endpoint
```warp-runnable-command
curl -s https://$HOST/api/v0/ping | jq -r .

```
This should return\:
```json
{
  "status": "ok"
}
```
## Activity data
### Submit of activity data
Multiple records can be posted at ounce\, but will be stored as individual records\.
For demo purposes\, we can set an invalid country code\.
```warp-runnable-command
curl -s -X POST https://$HOST/api/v0/str/activity-data \
--header "Authorization: Bearer $TOKEN" \
--header "Content-Type: application/json" \
--data '{
  "data": [
    {
      "numberOfGuests": 10,
      "countryOfGuests": [
        "ITZ",
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
    },
    {
      "numberOfGuests": 2,
      "countryOfGuests": [
        "BEL"
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
    },
    {
      "numberOfGuests": 3,
      "countryOfGuests": [
        "BEL"
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
' \
| jq .
```
Data has been pushed to the Kafka topic\: `activity-data`

### Only for demonstration purposes\, and if you have an own deployment 
We can use `kcat` to list kafka topics and kafka topics content
```warp-runnable-command
kcat -L
```
```warp-runnable-command
kcat -C -t activity-data -c 20 -f 'Headers: %h || Message value: %s\n'
```
Confluent Kafka also allows to views data in the topics\: [https\:\/\/confluent\.cloud\/environments\/env\-d19ry\/clusters\/lkc\-1dm6dj\/topics\/activity\-data\/message\-viewer](https://confluent.cloud/environments/env-d19ry/clusters/lkc-1dm6dj/topics/activity-data/message-viewer)
Note also the header data contains the OAUTH2 app name\.

### Retrieve activity data submitted to the SDEP
```warp-runnable-command
curl -s https://$HOST/api/v0/ca/activity-data \
--header "Authorization: Bearer $TOKEN" \
| jq .
```
Note each consumer has it\'s offset and will only get \"new\" published data\.

## Listings
### Submit of listings data
```warp-runnable-command
curl -s -X POST https://$HOST/api/v0/str/listings \
--header "Authorization: Bearer $TOKEN" \
--header "Content-Type: application/json" \
--data '
{
  "data": [
    {
      "registrationNumber": "1234",
      "Unit": {
        "description": "string",
        "floorLevel": "string",
        "address": {
          "street": "Culliganlaan 5",
          "city": "Diegem",
          "postalCode": "1831",
          "country": "BEL"
        },
        "obtainedAuth": true,
        "subjectToAuth": true,
        "numberOfRooms": 0,
        "occupancy": 0,
        "purpose": "string",
        "type": "string",
        "url": "STR-Platform.com/1234"
      }
    }
  ],
  "metadata": {
    "platform": "STR-Platform"
  }
}
'\
| jq .
```
### Only for demonstration purposes\, and if you have an own deployment
Consuming data from the kafka topic `listings`
```warp-runnable-command
kcat -C -t listings -c 20 -f 'Headers: %h || Message value: %s\n'
```
### Retrieving listing data
The limit parameter is optional\.
This data retrieval is also incremental and will only return \"new\" data\.
```warp-runnable-command
curl -s https://$HOST/api/v0/ca/listings \
--header "Authorization: Bearer $TOKEN" \
| jq .
```
## Area
### Shape file upload
```warp-runnable-command
curl -s -X POST https://$HOST/api/v0/ca/area \
--header "Authorization: Bearer $TOKEN" \
-F "file=@/Users/thierryturpin/GolandProjects/str-ap-internal/BELGIUM_-_Regions.shp"
```
### List shape files
```warp-runnable-command
curl -s https://$HOST/api/v0/str/area \
--header "Authorization: Bearer $TOKEN" \
| jq .
```
### Download of a shape file
Set the ID returned by the previous command
```warp-runnable-command
curl -s https://$HOST/api/v0/str/area/01J0R0GC700000000000000000 \
--header "Authorization: Bearer $TOKEN" \
-o downloaded_shape_file.shp
```
Verify the downloaded file
```warp-runnable-command
file downloaded_shape_file.shp

```
### Verify SSL certificate
Only relevant for own deployments\:
```warp-runnable-command
export SERVER_NAME=$HOST
export PORT=443 
openssl s_client -servername $SERVER_NAME -connect $SERVER_NAME:$PORT | openssl x509 -noout -dates
```
