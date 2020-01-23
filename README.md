# Wrike Timelog Extractor

A tool developed for utility reasons, making it easy to extract data from Wrike using its API in a format that can then be imported into another system (SAP in this case).

Built in Golang just for east of deployment (single file goodness)

Requires that a Wrike permanent access key (which can be created under API and Integrations, a menu option under a user profile) is passed via the command line or as an environment variable called 'WRIKEKEY'.