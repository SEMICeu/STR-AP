

**Disclaimer:**  
This document is an ongoing resource and may be updated regularly as new information becomes available or as further clarifications are provided. Please check back for the latest version to ensure you have the most current details.

Version: v01

# STR Q&A Webinar – 21/11/25

## Random Checks & Validation

**Q1:** We assume that the random checks conducted by booking platforms will be undertaken via the SDEP, as set out in the documentation on the Github. If so, is it designed so that the random checks will not only check if the registration number is valid, but that it is a valid registration number for that unit address and host?  
**A1:** The SDEP implementation is designed to ensure that random checks validate only the registration number and additionally gives back the address and name of the host. The SDEP would allow to cater for a certain number of information which the MS put at the disposal of the platforms, therefore if the MS put at the disposal of the platforms only the registration number the platforms will check only that the registration number indeed exist but if the MS put at the disposal of the platforms also the address and the name of the host then the platforms will check that all the elements correspond.

**Q2:** Have any parameters for the frequency of random checks by booking platforms been decided or agreed with the platforms?  
**A2:** No, this has not been addressed yet. There are currently no agreed parameters or guidelines on the frequency of random checks by booking platforms. This aspect remains open and will need to be defined.

**Q3:** Are there any guidelines to the frequency and selection size of random checks? Is any enforcement possible?  
**A3:** No formal agreement on this yet.

**Q4:** How can a municipality gain insight into the execution and outcomes of random checks done by platforms?  
**A4:** Article 7(2) of the Regulation provides that platforms shall inform the competent authorities of the results of the random checks.

**Q5:** What is the state of misleading self-declaration of hosts on platforms? Any advances on how to prevent it / spot such units during the random checks?  
**A5:** As indicated in article 7(1) lett. C) STR platforms have to make reasonable efforts to randomly check on a regular basis, declarations of the hosts concerning the existence or not of a registration procedure (the self-declaration).

---

## Data Submission & Deadlines

**Q6:** Regarding the submission of activity data by platforms, on the Github, the table under scenario A - user story 2 includes date timeframes, which we presume is to differentiate between stays and to avoid double counting of guest nationalities for example. We welcome this and we wish to clarify if the platforms have agreed to submit the data in this format?  
**A6:** The prototype and user stories describe the expected data model and endpoints for activity data submission (including rental period and guest details), and the table on GitHub appears to support differentiation between stays to avoid double counting.

**Q7:** Has a specific date or window been agreed for the monthly/quarterly submission of data from the platforms under Article 9? E.g. will it be at the end of each month/quarter? Is the expectation that platforms must submit data for April 2026 (or the previous quarter) by the end of May 2026, as the STR will be in effect at that point?  
**A7:** No, a specific date or submission window has not yet been agreed for monthly or quarterly data submissions under Article 9. There is currently no confirmed expectation for deadlines such as submitting April 2026 data by the end of May 2026. This point remains open and will be a point in the agenda of the next SDEP Coordination Group meeting on 18/12.

**Q8:** Should a submission deadline be set for STR platforms monthly activity data submissions? Should this be uniform for the Union? If so, what should the monthly deadline be?  
**A8:** No, a specific date or submission window has not yet been agreed for monthly or quarterly data submissions under Article 9. This point remains open and will be discussed during the next SDEP coordination group meeting on 18/12.

**Q9:** What is the status of the connection to the SDEP for platforms that facilitate short-stay rentals (1-6 months)? Is this process (data delivery) being evaluated? If so, when?  
**A9:** We will come back to this question. 

**Q10:** BE will add “purpose of stay” as (optional) field for activity data submissions. Are other member states considering the same?  
**A10:** This can be further discussed during the coordination group on December 18.

---

## Technical Specifications & APIs

**Q11:** Is there an OpenAPI (OAS) STR API specification available?  
**A11:** Swagger is available via the following link: https://eu-str.sdep-pilot.eu/swagger/index.html

**Q12:** Please confirm the API's functionality for verifying the STR object's registration code. In accordance with the regulation, can platforms be required to verify the validity of the object's registration code using the API service each time a client clicks the reservation button?  
**A12:** The platforms will not accept a technical dependency on an external API, which also reduces the SLA that would otherwise be required by the SDEP. So the validity is not part of the platforms validation process. It should also be mentioned that, according to the DMA, validation may not be legally required for every registration number

**Q13:** Has any further clarity been obtained regarding the updating by booking platforms of listings to include the field for registration number from 20th May 2026? Is it known if platforms are going to roll out automated mandatory updates to hosts to require them to fill the registration number field?  
**A13:** Not known at this stage.

**Q14:** The documentation states that all exports must be available in JSON, CSV, and Excel formats. In the French application, in addition to multiple exports, it is possible to display and export data for each Registration Number individually. We believe that exporting in JSON format is not particularly useful in the specific case of exporting data for a single Registration Number. Can we omit the JSON export option in this case?  
**A14:** We will come back on this question. 

---

## Regulatory Compliance & Enforcement

**Q15:** What actions can we expect from PwC / EU commission regarding interoperability assessment (Interoperable Europe Act)?  
**A15:** The STR Regulation as such has been adopted in 2024, pre-dating the applicability of the Interoperable Europe Act and thus not falling under the interoperability assessment obligation. However, if the Single Digital Entry Point is further described (by means of requirements) in autonomous acts and these acts were adopted after 12 January 2025, they would eventually trigger the interoperability assessment.

**Q16:** Are there any changes in the regulation or in the agreements with rental platforms? If so, what are they?  
**A16:** We have not had further specific contacts with platforms but they are abide by the specifications defined.

**Q17:** How is enforcement handled if platforms fail to meet the obligation to share data?  
**A17:** MS adopt rules on sanctions which will be imposed on platforms in case they fail to meet obligations.

**Q18:** What happens if rental platforms do not provide the data or do so inconsistently?  
**A18:** Article 15(2) provides that is up to MS to provide adequate sanctions in these cases.

**Q19:** Can member states report if rental platforms do not comply with the rules/agreements? If so, where? Can this also be done via the SDEP?  
**A19:** This has not been discussed but this is the competence of MS to provide adequate sanctions. There is not a report mechanism because the responsibility lays within the MS.

**Q20:** Do the platforms themselves remove accounts/listings if a mandatory check determines that a host is violating the rules?  
**A20:** They are obliged to do so in case of actual knowledge (DSA art. 6).

**Q21:** How is it ensured that advertisers do not falsely declare that they are not a residential property and thus fall outside the STR rules? How can a government agency monitor this?  
**A21:** This is the issue we have discussed at length and the option which looks more effective is the combined trigger of STR Regulation and DSA together.

**Q22:** The limit for small/micro STR platforms is defined at union level (<4250 listings monthly average across the Union). How will this be monitored? Should platforms self-monitor?  
**A22:** We will come back to this question. 

---

## Platform Operations & Data Handling

**Q23:** Can you see in the SDEP from which platform the rental data originates?  
**A23:** This is visible in the metadata.

**Q24:** Is it true that data disappears from the SDEP after someone takes it?  
**A24:** Yes, the SDEP is only a transmission channel and does not store activity data permanently. Once the data is transmitted to the competent authority, it is no longer kept on the SDEP; retention and deletion are handled by the competent authority according to legal requirements.

**Q25:** How is it ensured that the data is available to both Eurostat and a municipality?  
**A25:** Member States must designate a national entity responsible for transmitting, for each unit, the activity data and registration numbers to national and regional statistical offices and to Eurostat on a monthly basis. The SDEP transmits activity data, registration numbers, and listing URLs to the competent authority for their territory.

**Q26:** What engagement has taken place at EU-level with smaller platforms who will be subject to the requirements of the STR? How is it intended that they will submit their data manually via the SDEP, or is this for each Member State to decide at national level?  
**A26:** The technical specifications indicate that the SDEP will support simplified data submissions and interoperability with third-party systems, enabling manual uploads for smaller platforms. MS remain free to agree the way to submit the data.

**Q27:** What experience is there with using the prototype, including the extent to which it is assessed to reduce the costs of implementation?  
**A27:** We will come back to this question. 

---

## Registration Number Format & Retention

**Q28:** Could you please inform me whether the European Commission intends to adopt an implementing act establishing a harmonised format for the registration number?  
**A28:** There will not be, at least for the time being, any implementing act. The harmonised format of the registration number is already specified in the technical specifications prepared by PwC.

**Q29:** When not specified by STR regulations, how long is data retained once a Registration Number is deactivated? Is the retention period still 18 months, by analogy with cases where the host deletes their own Registration Number and with the retention period for activity data?  
**A29:** According to the GDPR personal data can be kept no longer than is necessary for the specific purpose. So we consider in that case the retention of the registration number in case of deactivation should be the same as in case of deletion by the host.

---

## Other Implementation Issues

**Q30:** Which countries / how many countries are set to implement STR?  
**A30:** This is one of the points in the agenda of next SDEP Coordination Group meeting on 18/12.

**Q31:** What are the actions that a MS which opts-out has to implement for the Regulation to work in the EU / intranational level? (e.g. a platform is based in a MS which opted-out)  
**A31:** No actions are needed from a MS which do not want data from platforms.

**Q32:** STR-AP unit/number of bedrooms > as confirmed by Paolo, based on the definition, this should be number of bed spaces > can the attribute title be changed?  
**A32:** This change can be submitted via a change request.
