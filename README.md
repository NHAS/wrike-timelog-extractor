﻿# Wrike Timelog Extractor

A tool developed for utility reasons, making it easy to extract data from Wrike using its API in a format that can then be imported into another system (SAP in this case).

Built in Golang just for east of deployment (single file goodness)

Requires that a Wrike permanent access key (which can be created under API and Integrations, a menu option under a user profile) is passed via the command line or as an environment variable called 'WRIKEKEY'.

For debugging, the launch.json (used with VS Code) refers to a local environment variable file called '.env'. If you create this, add the environment variable as WRIKEKEY="[key here]" into the file, you can debug with the key set. .env is part of .gitignore so no keys are leaked this way.

## Wrike API docs

This go prog calls the wrike apis as mentioned, and the documentation for these are [here](https://developers.wrike.com/documentation/api/overview). Generally they are fairly easy to use, and the code includes some special use like passing parameters (including passing json arrays in the query string, which is odd).
