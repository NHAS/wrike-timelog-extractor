package main

import (
	"encoding/json"
	"models"
)

func getContactsMap(apiKey string) (result map[string]Contact, err error) {

	textContent := getDataForURL(host+"/contacts", apiKey)

	var contacts Contacts
	err = json.Unmarshal(textContent, &contacts)
	if err != nil {
		return result, err
	}

	result = make(map[string]Contact)

	for _, field := range contacts.Data {

		result[field.Id] = field
	}

	return result, nil
}

func getFoldersAsTasks(apiKey string) (result map[string]collectiveTimeLog, err error) {
	textContent := getDataForURL(host+"/folders", apiKey)

	var folders collectiveTimelogs
	err = json.Unmarshal(textContent, &folders)
	if err != nil {
		return result, err
	}

	// the root folder api does not return custom fields
	// so need to compose ids, then call the folder api with these
	ids := ""
	for _, k := range folders.Data {
		ids += k.Id + ","
	}

	textContent = getDataForURL(host+"/folders/"+ids, apiKey)

	err = json.Unmarshal(textContent, &folders)
	if err != nil {
		return result, err
	}

	result = make(map[string]collectiveTimeLog, 0)

	for _, k := range folders.Data {
		if len(k.CustomFields) != 0 {
			result[k.Id] = k
		}
	}

	return result, nil
}

func getCustomFieldsMap(apiKey string) (result map[string]CustomField, err error) {
	textContent := getDataForURL(host+"/customfields", apiKey)

	var fields CustomFields
	err = json.Unmarshal(textContent, &fields)
	if err != nil {
		return result, err
	}

	result = make(map[string]CustomField)

	for _, field := range fields.Data {
		result[field.Id] = field
	}

	return result, nil
}

func getTimelogMap(apiKey string, contacts map[string]Contact) (result map[string][]Timelog, err error) {
	textContent := getDataForURL(host+"/timelogs", apiKey)

	var timelogs Timelogs
	err = json.Unmarshal(textContent, &timelogs)
	if err != nil {
		return result, err
	}

	result = make(map[string][]Timelog)
	for _, k := range timelogs.Data {

		k.User = contacts[k.UserId].FirstName + " " + contacts[k.UserId].LastName
		result[k.TaskId] = append(result[k.TaskId], k)
	}

	return result, nil
}
