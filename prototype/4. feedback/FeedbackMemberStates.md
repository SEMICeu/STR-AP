| Disclaimer  |
|-----------|
| This report was prepared for DG Grow by PwC EU Services. The views expressed in this report are purely those of the authors and may not, in any circumstances, be interpreted as stating an official position of the European Commission. The European Commission does not guarantee the accuracy of the information included in this report, nor does it accept any responsibility for any use thereof. Reference herein to any specific products, specifications, process, or service by trade name, trademark, manufacturer, or otherwise, does not necessarily constitute or imply its endorsement, recommendation, or favouring by the European Commission. All care has been taken by the author to ensure that s/he has obtained, where necessary, permission to use any parts of manuscripts including illustrations, maps, and graphs, on which intellectual property rights already exist from the titular holder(s) of such rights or from her/his or their legal representative.|

# 1. Introduction  
The survey witnessed the participation of 10 member states, namely Belgium, France, Germany, Hungary, Ireland, Luxembourg, Netherlands, Portugal, Spain, and Austria. 

# 2. Feedback Member States

## 2.1. API Architecture 

## 2.1.1 General Results 

All Member States that participated in the STR survey have reached a consensus on using a RESTful API architecture for the Single Digital Entry Point (SDEP). This agreement signifies that there is a unified understanding and acceptance of the chosen API architecture among the participating Member States.

*CONSENSUS - YES: 10/10*

**Way forward - Prototype will be using a RESTful API**

## 2.1.2 Questions/Remarks Member States

“Yes, we support the use of a RESTful API and suggest going further down this path and specifying the API based on the OpenAPI 3.0 specification and the JSON schema rather than using JSON-LD, not least because OpenAPI based on JSON schema is more widely used and offers much better tool support.”

## 2.2. Security and Authentication

## 2.2.1 General Results 
During the survey, it was observed that there was a remarkable level of consensus among the participating member states regarding the security protocol OAuth 2.0. OAuth 2.0 is widely recognized as a robust and effective framework for securing access to APIs and protecting user data. The majority of the 10 member states expressed their support for OAuth 2.0 as a reliable and standardised solution for authentication and authorization.

*CLOSE TO CONSENSUS - YES: 8/10*

**Way forward - Prototype will be using OAuth 2.0 for Authentication.**

## 2.2.2 Questions/Remarks Member States
“We prefer using API keys and IP filtering. We would rather use this kind of security as it seems efficient enough and less costly.”

“OIDC, which extends Oauth2.0.”

## 2.3. Data Format

## 2.3.1 General Results
In the survey, there was a strong consensus among the participating member states regarding the use of JSON (JavaScript Object Notation) as a data exchange format. JSON has gained widespread popularity and acceptance due to its simplicity, flexibility, and compatibility with various programming languages and platforms. The majority of the 10 member states expressed their support for JSON as a preferred format for exchanging data between systems and applications. JSON's lightweight structure and human-readable syntax make it easy to understand and work with, facilitating efficient data transmission and interoperability.

*CLOSE TO CONSENSUS - YES: 9/10*

**Way forward - Prototype will be using a JSON format with the exception of the shapefile/area data exchange.**

## 2.3.2 Questions/Remarks Member States
“We may receive a big amount of data from STR platforms. The CSV format is lighter than JSON and may facilitate the management of resources.”

## 2.4. SDEP Endpoints

## 2.4.1 GET/Area: General Results
Several Member States have shared their feedback on different endpoints of the SDEP. The majority of respondents do not support the GET:/area endpoint. 

*NO CONSENSUS - YES: 5/10 (Name and purpose)*

**Way forward - Prototype will investigate a solution where for MS that have multiple CA’s to work with shapefiles and how we can have a system that the responsibility for uploading the area’s is in the hands of CA. For MS that have a national regulation a more suited approach will be investigated.**

## 2.4.2 GET/Area: Questions/Remarks Member States
“We consider allowing STR platforms to download a file containing city names and insee codes (that are more precise than postal codes). It seems easier to manage and to send to STR platforms. ”

“We recommend implementing an event-log concept comparable to an Apache Kafka "compacted" idempotent topic with key. On the one hand, this should prevent the complete inventory data from having to be delivered with every GET, even though it is 99% unchanged. On the other hand, the query is decoupled from the provision of the data by the responsible authorities.
Specifically, the following query parameters should also be included in the GET:
offset: Specifies the position in the event log from which the data should be returned. An offset of 0 indicates that all aeras are to be delivered in full
pagination: Number of data records per response
country: ISO 3166-1 code of the country to be queried (optional)
area: Area identifier to restrict the query to a specific area (optional)”

## 2.4.3 POST/Listings: General Results
Several Member States have shared their feedback on different endpoints of the SDEP. The majority of respondents do not support the POST:/Listings endpoint. 

*NO CONSENSUS - YES: 4/10 (Name and purpose)*

**Way forward - Only the data exchange based on the agreed datamodel will be part of the prototype. A way to align on the functional requirement between MS and Platforms will be discussed but is outside of scope of the Prototype.**

## 2.4.4 POST/Listings: Questions/Remarks Member States

## 2.4.5 GET/Orders: General Results
Several Member States have shared their feedback on different endpoints of the SDEP. The majority of respondents do not support the POST:/Listings endpoint. 

*NO CONSENSUS - YES: 5/10 (Name and purpose)*

**Way forward - Only the data exchange based on the agreed datamodel will be part of the prototype. A way to align on the functional requirement between MS and Platforms will be discussed but is outside of scope of the Prototype.**

## 2.4.6 GET/Orders: Questions/Remarks Member States

"We recommend implementing an bidirectional event-log concept comparable to an Apache Kafka "compacted" idempotent topic with key.
On the one hand, this should prevent the complete inventory data from having to be delivered with every GET, even though it is 99% unchanged. On the other hand, the query is decoupled from the provision of the data by the responsible authorities or online platforms.
 
Benefits:
Platforms no longer need to submit their entire inventory of listings. Instead, only those listings that have undergone a change will be transmitted. In addition, a flag must be introduced to indicate that a listing may be suspended or deleted.
Competent authorities may send more than one response. If the order is in "pending" status, i.e. a registration number has been assigned but the registration process has not yet been completed. Once completed, the status can be sent as "approved".
This asynchronous update becomes even more interesting if the online platform has previously been informed that the competent authority intends to suspend a certain unit. If the host subsequently provides the missing information, the Competent Authority can automatically communicate the changed status back to the online platforms involved.
 
Specifically, the following query parameters should also be included in the GET:
offset: Specifies the position in the event log from which the data should be returned. An offset of 0 indicates that all aeras are to be delivered in full
pagination: Number of data records per response
country: ISO 3166-1 code of the country to be queried (optional)
registrationNumber: To be able to limit the search to a specific order (optional)"

## 2.4.7 POST/Activity-data: General Results
Almost All Member States agreed on the name and purpose of this endpoint

*CONSENSUS - YES: 9/10 (Name and purpose)*

**Way forward - SDEP endpoint to be built as collaborative as possible.**

## 2.4.8 POST/Activity-data: Questions/Remarks Member States

## 2.5. Response Codes

## 2.5.1 General Results
Almost all Member States agreed on the recommended response codes.

*CONSENSUS - YES: 9/10*

**Way forward - Agreed response codes will be used when developing the endpoints.**

## 2.5.2 Questions/Remarks Member States
/

# Questions

**Appeals mechanism: the Regulation states, at Article 4(3)(c), that registration procedures are subject to effective appeal mechanisms within the Member State. Would the Commission/consultants have any more information on how this provision is intended to operate, or is it for each Member State to determine the manner in which such an appeal mechanism is given effect?**

The Regulation states in Article 4(3)(c) that registration procedures are subject to effective appeal mechanisms within the Member State. However, the document does not provide further information on how this provision is intended to operate. 

**Article 4(5) states that 'Member States shall ensure that registration numbers are included in a public and easily accessible registry. As part of our registration system, we intend to offer a search box function which would allow a member of the public, or any other user, to check the validity of a registration number. Would this meet the requirements of Article 4(5)?**

As part of your registration system, you intend to offer a search box function to check the validity of a registration number. This would likely meet the requirements of Article 4(5) as it would provide a means for the public or any other user to easily access and verify the validity of a registration number.

**During the 7th March webinar we heard different information regarding the /listings API. One person said the platforms will send data to this API if a property does not have a valid registration number. Another said the platforms would send data to this API "regularly" regardless. We wish to seek clarity on this matter.**

The sending of data for a non-valid registration number would likely be a result of random checks.  On the other hand, the regular sending of data would pertain to the activity data sent by the platforms. 

**Is activity reporting done per listing or per registration number?**

The regulation does not explicitly state whether activity reporting is done per listing or per registration number

**If per listing: is the listing url the unique identifier or is an additional identifier applicable?**

We foresee the url as a unique identifier. 

**Include the registration numbers within the activity report object**

Should indeed be added. 

**Discuss how to handle aggregrate data for statistical purposes (e.g. Eurostat, CBS)**

Not part of the prototype

**Discuss concrete suggestions how access control is organised (e.g. distribution of API-keys etc)**

Not part of the prototype

**We notice that the all flows assume the happy flow, no error handling is in place yet and should be documented more concretely.**

Although this will not be part of the prototype phase, this indeed should be discussed in more detail. 

**How are listing that do not require a registration number (e.g. hotels) handled.**

More discussion is needed on the definitions of the unit types. We have not found this to be in consensus. We propose following a standard definition based on schema.org. 

**How is data shared by smaller platforms?**

The data exchange will be done in the same way regardless of the size. The exception of the very small platforms that can provide the activity data in an alternative way if they find this better suited. 
