
# EU STR - Single Digital Entry Point
| Disclaimer  |
|-----------|
| This report was prepared for DG Grow by PwC EU Services. The views expressed in this report are purely those of the authors and may not, in any circumstances, be interpreted as stating an official position of the European Commission. The European Commission does not guarantee the accuracy of the information included in this report, nor does it accept any responsibility for any use thereof. Reference herein to any specific products, specifications, process, or service by trade name, trademark, manufacturer, or otherwise, does not necessarily constitute or imply its endorsement, recommendation, or favouring by the European Commission. All care has been taken by the author to ensure that s/he has obtained, where necessary, permission to use any parts of manuscripts including illustrations, maps, and graphs, on which intellectual property rights already exist from the titular holder(s) of such rights or from her/his or their legal representative.|

## Introduction to the EU STR - Single Digital Entry Point Prototype
This report provides a guide on setting up the SDEP Prototype.  The goal is to enable smooth and interoperable information among Member States and STR platforms, with efficient access, data integrity, and security.

### Purpose and Scope

The primary aim of the EU STR - Single Digital Entry Point prototype is to create a secure channel for the electronic transmission of data related to short-term rentals. The prototype prioritizes the development and integration of three key data categories:

1. **Listings**
2. **Activity Data**
3. **Area**

### Features and Capabilities

The prototype offers several robust features designed to streamline data exchange:

- **Secure Data Transmission**: Utilizes OAuth 2.0 for secure authentication and authorization.
- **Comprehensive API Endpoints**: Provides endpoints for submitting and retrieving listings, activity data, and geospatial area data. 
- **Data Validation**: Ensures the integrity and accuracy of the transmitted data.
- **Incremental Data Retrieval**: Consumers receive only the new data, optimizing the data handling process.
- **Kafka Integration**: Utilizes Kafka for efficient data streaming and processing.

### Audience

This prototype is designed for testing purpose for:

- **Short-Term Rental Platforms**
- **Competent Authorities**: 

### Getting Started

To begin using the EU STR - Single Digital Entry Point API, users need to:

1. Obtain API credentials by contacting the designated support email.  
2. Set up their environment by configuring client credentials and obtaining an OAuth token.
3. Use the provided API endpoints to submit and retrieve data as per their requirements.

The following documentation is divided into 3 parts: 
- **1. Technical: EU STR - Code Single Digital Entry Point**
- **2. Technical: EU STR - Infra Single Digital Entry Point**
- **3. Testing¨**

The first 2 folders will guide you in setting up the SDEP prototype yourself. The last folder will help you understanding the prototype, but also will help you in testing the prototype (that is hosted on our servers). 

# 1. Technical: EU STR - Single Digital Entry Point
## Purpose of This folder
This folder serves as a practical guide for developers, IT teams, and other stakeholders to validate and interact with the SDEP API. It includes:

Setup Instructions: Steps to configure your (own) environment, obtain OAuth tokens, and set API host details. However, with the base url, already the api can be tested (hosted on our servers). (=> See also folder Testing)
API Endpoints: curl commands to test various API endpoints, including health checks, submission of activity data, listings data, and area data (shapefiles)
Data Formatting: Usage of jq for formatting JSON responses for readability.
Kafka Integration: Demonstration of data handling using Kafka topics for real-time data streaming and processing.

## Prerequisites
To use this notebook effectively, you will need:

Bitwarden CLI: For managing client secrets securely.

curl and jq: Installed on your terminal for making API requests and formatting JSON responses.

Kafka Tools (Optional): For viewing and managing Kafka topics if you have your own deployment.

## How to Use This Notebook
Configure Client Credentials: Replace placeholders with your actual CLIENT_ID and CLIENT_SECRET.
Obtain OAuth Token: Follow the provided steps to get an access token for authenticating API requests.
Test API Endpoints: Execute the curl commands to interact with the endpoints and validate the API functionality.
Monitor Data Streams (Optional): If you have a Kafka deployment, use the provided commands to monitor data topics.

[Link to notebook](https://github.com/SEMICeu/STR-AP/tree/main/prototype/1.%20Technical_%20EU%20STR%20-%20Code%20Single%20Digital%20Entry%20Point)


# 2. Technical: EU STR - Infra Single Digital Entry Point**
## Purpose of This folder
The primary goal of this folder is to guide developers and infrastructure engineers through the process of deploying a new version of the SDEP on an EKS cluster. The deployment process involves building a Docker image, updating Helm Charts, and using Pulumi to manage deployments. This notebook also covers the configuration of Kafka topics, AWS infrastructure, and Helm Chart deployments. 

## Prerequisites
To effectively use this notebook, you will need:

- **Pulumi CLI**: For managing infrastructure as code.
- **Kubectl**: For interacting with the Kubernetes cluster.
- **AWS CLI**: For AWS account configuration and management.
- **Docker**: For building application images.
- **Access to Repositories**: Ensure you have access to the relevant GitHub repositories and Docker registries.

## How to Use This Notebook
This folder  is a comprehensive guide for managing and deploying infrastructure and applications using Pulumi. It includes detailed instructions for configuring client credentials, setting up Kafka and AWS configurations, and deploying Helm Charts on an EKS cluster. The notebook covers the entire workflow, from building Docker images and updating Helm Charts to verifying deployments using kubectl. Additionally, it provides steps for bootstrapping an AWS account, including activating MFA and configuring IAM users and roles. 

[Link to notebook](https://github.com/SEMICeu/STR-AP/tree/main/prototype/2.%20Technical_%20EU%20STR%20-%20Infra%20Single%20Digital%20Entry%20Point)

# 3. Testing
## Purpose of This folder
The primary goal of this folder is to guide developers through the process of testing the SDEP prototype that we hosted on our own environment. This folder also includes dummy data, API endpoints documentation, a Postman collections for automated testing, and documentation to understand the technical architecture.

**API Credentials: Obtain these by contacting this support email: wouter.travers@pwc.com and victor.vanhullebusch@pwc.com**

## Prerequisites
This folder includes resources to help you in testing the prototype.  It includes detailed instructions for: 

- **Postman**: For API testing and automation.
- **Access to Dummy Data**: Sample data to facilitate testing scenarios.
- **Access to the Slide Deck**: For understanding the technical components of the prototype.
- **Knowledge of API Endpoints**: Familiarity with the API endpoints.

[Link to notebook](https://github.com/SEMICeu/STR-AP/tree/main/prototype/3.%20Testing)

# 4. Feedback 
[Link to notebook](https://github.com/SEMICeu/STR-AP/tree/main/prototype/4.%20Feedback%20)



