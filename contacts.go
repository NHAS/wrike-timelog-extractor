package main

import (
	"encoding/json"
)

type Contact struct {
	Id        string
	FirstName string
	LastName  string
}

type Contacts struct {
	Data []Contact
}

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
