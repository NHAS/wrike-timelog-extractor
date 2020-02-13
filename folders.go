package main

import (
	"encoding/json"
)

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
