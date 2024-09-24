
<p align="center">
  <img src="images/STRHeader.png" alt="STR framework">
</p>

| <span style="background-color:#1c4fa4; color:white">Settings</span>  | Value |
|-------------------|----------------------|
| Document Title   | D06.01 Report on the Prototype –Platforms |
| Project Title | Interoperability solutions in the area of Short-Term Rental (STR) services |
| Document Author | PwC EU Services |
| Project Owner |  DG GROW - European Commission |
| Project Manager | Travers Wouter  - PwC EU Services |

| Disclaimer  |
|-----------|
| This report was prepared for DG Grow  by PwC EU Services. The views expressed in this report are purely those of the authors and may not, in any circumstances, be interpreted as stating an official position of the European Commission. The European Commission does not guarantee the accuracy of the information included in this report, nor does it accept any responsibility for any use thereof. Reference herein to any specific products, specifications, process, or service by trade name, trademark, manufacturer, or otherwise, does not necessarily constitute or imply its endorsement, recommendation, or favouring by the European Commission. All care has been taken by the author to ensure that s/he has obtained, where necessary, permission to use any parts of manuscripts including illustrations, maps, and graphs, on which intellectual property rights already exist from the titular holder(s) of such rights or from her/his or their legal representative.|

# Table of Content

1. [Executive Summary](#1-executive-summary)  
2. [Introduction](#2-introduction)  
3. [Overview of the EU Regulation and Requirements](#3-overview-of-the-eu-regulation-and-requirements)  
   - 3.1. [Background and objectives](#31-background-and-objectives)  
   - 3.2. [Key compliance aspect](#32-key-compliance-aspect)  
4. [API Endpoints and Usage](#4-api-endpoints-and-usage)  
   - 4.1. [Endpoint specifications](#41-endpoint-specifications)  
      - 4.1.1. [General Endpoints](#411-general-endpoints)  
      - 4.1.2. [Endpoints for Platforms](#412-endpoints-for-platforms)  
   - 4.2. [Authentication and authorization](#42-authentication-and-authorization)  
      - [OAuth 2.0 Framework Roles](#oauth-20-framework-roles)  
5. [Technical Infrastructure](#5-technical-infrastructure)  
   - [User and Server Requests](#user-and-server-requests)  
   - [Network Load Balancer (NLB)](#network-load-balancer-nlb)  
   - [Nginx Ingress Controller](#nginx-ingress-controller)  
   - [Kubernetes Service](#kubernetes-service)  
   - [Pods](#pods)  
   - [Persistent Volume Claim (PVC)](#persistent-volume-claim-pvc)  
   - [Apache Kafka Integration](#apache-kafka-integration)  
   - [Infrastructure Management](#infrastructure-management)  
6. [Testing Steps](#6-testing-steps)  
   - 6.1. [Via Terminal Commands](#61-via-terminal-commands)  
      - 6.1.1. [Authentication](#611-authentication)  
      - 6.1.2. [Get the OAUTH token (from the /token endpoint)](#612-get-the-oauth-token-from-the-token-endpoint)  
      - 6.1.3. [Define the HOST](#613-define-the-host)  
      - 6.1.4. [Health check endpoint test (endpoint 1 for platforms)](#614-health-check-endpoint-test-endpoint-1-for-platforms)  
      - 6.1.5. [Submitting activity data endpoint (endpoint 2 for platforms)](#615-submitting-activity-data-endpoint-endpoint-2-for-platforms)  
      - 6.1.6. [Download Shapefiles Uploaded by Competent Authorities (endpoint 3 for platforms)](#616-download-shapefiles-uploaded-by-competent-authorities-endpoint-3-for-platforms)  
      - 6.1.7. [Download List of Uploaded Shapefiles by Competent Authorities (endpoint 4 for platforms)](#617-download-list-of-uploaded-shapefiles-by-competent-authorities-endpoint-4-for-platforms)  
   - 6.2. [Via Postman](#62-via-postman)  


# 1. Executive Summary  
  
The D06.01 Report on the Prototype – Platforms – provides a comprehensive overview of the interoperability solutions in the area of Short-Term Rental (STR) services, developed in response to the STR European Union regulations. The report, prepared by PwC EU Services for DG GROW - European Commission, outlines a prototype with best practices for the technical and regulatory framework necessary for STR platforms.  
  
The report begins with an introduction to the EU regulation mandating the reporting of short-term rental activities. It details the regulation's background, objectives, and key compliance aspects, emphasizing the necessity for a Single Digital Entry Point (SDEP) to facilitate data integration and reporting.  
  
The document provides detailed specifications for the API endpoints developed for data transmission. It includes practical examples of GET and POST requests, guiding developers on how to interact with the SDEP effectively and securely. The endpoints cover health checks, activity data submission, shapefile downloads, and invalid listing reports.  
  
The report delves into the technical backbone of the prototype, explaining the use of Kubernetes for container orchestration, Pulumi for infrastructure as code, Go for backend services, AWS for cloud solutions, Kafka for data streaming, and Helm charts for managing Kubernetes applications. This section highlights the robust and scalable architecture designed to ensure secure and efficient data transmission.  
  
Step-by-step instructions for setting up the environment and deploying the prototype are provided. This includes guidance on connecting to the API endpoints, ensuring accurate data transmission and validation, and leveraging tools like Postman for testing.  
  
By following the guidelines and leveraging the resources provided in this report, developers and integrators will be well-equipped to ensure compliance with the EU regulation, contributing to the overarching goal of transparent and accountable short-term rental activities. The prototype's sophisticated technology stack and contemporary architectural patterns guarantee secure, scalable, and efficient data transmission, aligning with EU regulatory requirements and enhancing overall efficiency.  

# 2. Introduction  
  
In response to the recent European Union regulation requiring short-term rental (STR) platforms to transmit activity data to designated authorities via a Single Digital Entry Point (SDEP), this document offers detailed guidance on connecting to the newly established API endpoints. These endpoints are integral to a prototype designed to ensure compliance with this regulation, facilitating smooth data integration and reporting.  
  
The prototype employs a sophisticated technology stack and contemporary architectural patterns to guarantee secure, scalable, and efficient data transmission. The employed technologies and methodologies include Kubernetes clusters for container orchestration, Pulumi for infrastructure as code, Go for backend service development, AWS Services for cloud infrastructure, Kafka architecture for reliable data streaming, Helm charts for Kubernetes application management, and standard API GET/POST requests for data interaction.  
  
This document is organized to provide a thorough overview of each component, offering explicit instructions and best practices for effective API endpoint usage. By adhering to the guidelines presented, developers and integrators will be well-equipped to ensure compliance with the EU regulation, contributing to the overarching goal of transparent and accountable short-term rental activities.  
  
By following the detailed instructions and leveraging the resources provided, stakeholders can achieve seamless integration with the Single Digital Entry Point, thereby aligning with EU regulatory requirements and enhancing overall efficiency.  
