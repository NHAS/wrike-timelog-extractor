package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getDataForURL(url string, apiKey string) map[string]interface{} {
	var bearer = "Bearer " + apiKey
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		os.Exit(2)
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var f interface{}
	json.Unmarshal([]byte(body), &f)

	return f.(map[string]interface{})
}

func getCustomFieldsMap(json map[string]interface{}) map[string]string {
	var data = json["data"].([]interface{})

	var result = make(map[string]string)
	for _, k := range data {
		var field = k.(map[string]interface{})
		result[field["id"].(string)] = field["title"].(string)
	}

	return result
}

func getContactsMap(json map[string]interface{}) map[string]string {
	var data = json["data"].([]interface{})

	var result = make(map[string]string)
	for _, k := range data {
		var field = k.(map[string]interface{})
		result[field["id"].(string)] = field["firstName"].(string) + " " + field["lastName"].(string)
	}

	return result
}

type timelog struct {
	user, trackedDate string
	hours             float64
}

func getTimelogMap(json map[string]interface{}, contacts map[string]string) map[string][]timelog {
	var data = json["data"].([]interface{})

	var result = make(map[string][]timelog)
	for _, k := range data {
		var entry = k.(map[string]interface{})
		var timelog = timelog{}

		var userID = entry["userId"].(string)
		timelog.user = contacts[userID]
		timelog.trackedDate = entry["trackedDate"].(string)
		timelog.hours = entry["hours"].(float64)

		var taskID = entry["taskId"].(string)
		result[taskID] = append(result[taskID], timelog)
	}

	return result
}

type task struct {
	fields   map[string]string
	timelogs []timelog
}

func getTasks(json map[string]interface{}, customFields map[string]string, timelogs map[string][]timelog) []task {
	var data = json["data"].([]interface{})

	var result = []task{}
	for _, k := range data {
		var entry = k.(map[string]interface{})
		var task = task{}

		task.fields = make(map[string]string)
		var entryFields = entry["customFields"].([]interface{})
		for _, j := range entryFields {
			var field = j.(map[string]interface{})
			var fieldName = customFields[field["id"].(string)]
			task.fields[fieldName] = field["value"].(string)
		}

		var taskID = entry["id"].(string)
		task.timelogs = timelogs[taskID]

		if len(task.fields) > 0 && len(task.timelogs) > 0 {
			result = append(result, task)
		}
	}

	return result
}

func main() {
	fmt.Println("wrike-extractor running!")

	var apiKey = os.Getenv("WRIKEKEY")
	if len(os.Args) == 2 {
		apiKey = os.Args[1]
	}

	if apiKey == "" {
		fmt.Println("please provide an API key (permananent access token) as an argument or as the env var WRIKEKEY")
		os.Exit(1)
	}

	var host = "https://www.wrike.com/api/v4"

	var json = getDataForURL(host+"/customfields", apiKey)
	var customFields = getCustomFieldsMap(json)

	json = getDataForURL(host+"/contacts", apiKey)
	var contacts = getContactsMap(json)

	json = getDataForURL(host+"/timelogs", apiKey)
	var timeLogs = getTimelogMap(json, contacts)

	json = getDataForURL(host+"/folders/IEADDL6NI4NMQXDT/tasks?fields=['customFields']", apiKey)
	var tasks = getTasks(json, customFields, timeLogs)

	fmt.Println(customFields)
	fmt.Println(contacts)
	fmt.Println(timeLogs)
	fmt.Println(tasks)

	// TODO: compose CSV and print out
}
