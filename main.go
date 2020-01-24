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

type taskTimeLog struct {
	fields  map[string]string
	timelog timelog
}

func getTaskTimelogs(json map[string]interface{}, customFields map[string]string, timelogs map[string][]timelog) []taskTimeLog {
	var data = json["data"].([]interface{})

	var result = []taskTimeLog{}
	for _, k := range data {
		var entry = k.(map[string]interface{})

		var fields = make(map[string]string)
		var entryFields = entry["customFields"].([]interface{})
		for _, j := range entryFields {
			var field = j.(map[string]interface{})
			var fieldName = customFields[field["id"].(string)]
			var fieldValue = field["value"].(string)
			if fieldValue != "" {
				fields[fieldName] = fieldValue
			}
		}

		if len(fields) > 0 {
			for key := range customFields {
				if _, exists := fields[key]; !exists {
					fields[key] = ""
				}
			}

			var taskID = entry["id"].(string)
			var taskLogs = timelogs[taskID]
			for _, log := range taskLogs {
				var task = taskTimeLog{}
				task.fields = fields
				task.timelog = log
				result = append(result, task)
			}
		}
	}

	return result
}

func main() {
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

	json = getDataForURL(host+"/tasks?subTasks=true&fields=['customFields']", apiKey)
	var tasks = getTaskTimelogs(json, customFields, timeLogs)

	var csv = ""
	for _, task := range tasks {
		for _, key := range customFields {
			csv += task.fields[key] + ","
		}
		csv += fmt.Sprintf("%s,%s,%f\n", task.timelog.user, task.timelog.trackedDate, task.timelog.hours)
	}

	fmt.Print(csv)
}
