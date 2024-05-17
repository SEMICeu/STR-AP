# User Story 1.1 - Host data submission

*As a host I want to submit my personal and unit information through the platform of the Competent Authority, so I can register my property for short term rental.*  

**Acceptance Criteria**
1. The host must be able to enter and update personal information online (and offline), including Name, ID number, Address, Contact Tel no, email, type of Unit, primary or secondary residence (or other), maximum no. of bed spaces/ guests through a user-friendly interface.
2. The system should allow the host to enter and modify unit details such as location, type, and other relevant attributes at any time.
3. The system must validate and confirm the registration number's issuance through an automated check, ensuring its uniqueness and the correctness of the associated host information.
4. Upon successful submission of the host information, the system must display a confirmation message to the host.

**Definition of Done**
1. The Host Information Submission Form is fully implemented, functional, and user- tested to confirm usability.
2. All data validation rules have been enforced on the client side, and appropriate error messages are reliably displayed for any incorrect or incomplete inputs.
3. The back-end API for storing host data has been developed, deployed, and tested to confirm data integrity and security.
4. Comprehensive integration tests have been completed and passed, confirming that the data submission process works end-to-end without any issues.

# User Story 1.2 - Registration Number Issuance

*As a host, I want to receive a unique registration number for each Unit I offer for STR which validates my property details so that I can officially list my property.*

**Acceptance Criteria:**

1. The system automatically and immediately issues a unique registration number upon successful data submission. Validation of data can happen afterwards.
2. Hosts can view the registration number in their account dashboard.
3. The registration number includes the ISO 3166-1 alpha-2 country code prefix and is unique per country.

**Definition of Done:**
1. Unique registration number generation logic is implemented and tested.
2. Hosts can view their registration number in their account dashboard.

# User Story 3.1: Activity Data Sharing

*As an STR platform, I want to share the required activity data with the SDEP.*

**Acceptance Criteria:**
1. The platform can compile and send activity data to the relevant stakeholders in a secure and compliant data format (e.g., API contract).
2. Activity data included in scope → see data model.
3. The system can automatically pull the latest data from the SDEP to the CA.
4. The Platforms can send data automatically to the SDEP in the agreed format. (see datamodel)

**Definition of Done:**
1. Activity data compilation and transmission facilitated as per developed compliance standards.
2. Data is sent (push) to SDEP endpoints ensuring data integrity.
3. The back-end API for storing host data has been developed, deployed, and tested to confirm data integrity and security.
4. Data is received (get) from SDEP endpoints for the Competent Authority.


# User Story 3.2: Area List Updates

*As a Competent Authority, I want to update the list of areas requiring registration so that the STR platform can check the veracity of the host declaration regarding the existence of registration in their area.*

**Acceptance Criteria:**
1. Competent Authorities can update the list of areas requiring registration through the system interface.
2. The system updates the area list dynamically, reflecting new data.
3. Member States are able to provide an API to call a Shapefile via the SDEP to the STR Platforms.
4. Optional: Member States are able to provide an API to give the list of postcodes or other data of the STR area's via the SDEP to the STR Platforms if the shapefile approach is not suited. To be discussed.

**Definition of Done:**
1. STR platform (or system) Lists of registration areas/ numbers dynamically updated upon receiving new data. The update is triggered manually and no reaction time is indicated. 

# User Story 3.3:  Dispatching Activity Data

*As a Competent Authority I want to get the activity data of the assigned STR Unit's in my Area.*

**Acceptance Criteria:**
1. SDEP dispatches data to the system for information to the Competent Authority.

**Definition of Done:**
1. An API contract is defined and agreed upon between the Competent Authority and the SDEP, that enables sharing Activity data. (See datamodel)
2. The functionality for the SDEP to ensure that only those Competent Authorities in charge of registration processes and the development and enforcement of local rules receive platform data.
3. The SDEP Platform's back-end is capable of receiving, acknowledging, processing, and responding to receiving activity data.

# User Story 4.1: Share Listings flagged in random checks

*As an STR platform, I want to comply with the regulation by performing random checks on listings and flag and communicate possible invalid listings to Competent Authorities via the SDEP.
Competent Authorities are not expected to respond directly to these requests. Competent Authorities will be able to respond via users story 4.2.*

**Acceptance Criteria:**
1. Listings with invalid or non-existing registration numbers flagged during the random check phase (to be defined), and the listings are sent to SDEP.
2. SDEP dispatches data to the Competent Authority.

**Definition of Done:**
1. An API contract is defined and agreed upon between the Competent Authority and the STR Platform, that enables sharingURL, registration number and property address. (See datamodel)
2. The functionality for the STR Platform to send potential invalid listings is implemented and tested, ensuring data is transmitted securely.
3. The SDEP Platform's back-end is capable of receiving, acknowledging and processing flagged listings.

# User Story 4.2: Order Removal Based on Listing URL and/or Registration Number

*As a Competent Authority I want to be able to send removal orders from STR Platforms based on Listing URL and /or registration number. This story does not discuss the how or the when, but strictly aims to showcase the possibility to take down orders as well as to define the API contract that can be used to align on the data exchange.*

**Acceptance Criteria:**
1. The Competent Authority can submit a removal order through an established interface or other electronic means, specifying the listing URL and/or registration number.
2. The STR Platform acknowledges receipt of the removal order.
3. The STR Platform processes the removal order and confirms the takedown of the specified listings.

**Definition of Done:**
1. An API contract is defined and agreed upon between the Competent Authority and the STR Platform, outlining the structure of removal orders and the expected responses.
2. The functionality for the Competent Authority to send removal orders is implemented and tested, ensuring data is transmitted securely.
3. The STR Platform's has a process capable of receiving, acknowledging, processing, and responding to removal orders.
