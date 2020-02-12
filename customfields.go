package main

import (
	"encoding/json"
)

type Field struct {
	Id    string
	Title string
}

type CustomFields struct {
	Data []Field
}

func getCustomFieldsMap(apiKey string) (result map[string]string, err error) {
	textContent := getDataForURL(host+"/customfields", apiKey)

	var fields CustomFields
	err = json.Unmarshal(textContent, &fields)
	if err != nil {
		return result, err
	}

	result = make(map[string]string)

	for _, field := range fields.Data {

		result[field.Id] = field.Title
	}

	return result, nil
}
