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
	resp, _ := client.Do(req)
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
	taskID, userID, trackedDate string
	hours                       float64
}

func getTimelogs(json map[string]interface{}) []timelog {
	var data = json["data"].([]interface{})

	var result = make([]timelog, len(data))
	for i, k := range data {
		var entry = k.(map[string]interface{})
		var timelog = timelog{}

		timelog.taskID = entry["taskId"].(string)
		timelog.userID = entry["userId"].(string)
		timelog.trackedDate = entry["trackedDate"].(string)
		timelog.hours = entry["hours"].(float64)

		result[i] = timelog
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
	var timeLogs = getTimelogs(json)

	// var tasks = getDataForURL(host+"/tasks?fields=['customFields']", apiKey)
	//

	fmt.Println(customFields)
	fmt.Println(contacts)
	fmt.Println(timeLogs)

	// TODO: compose CSV and print out
}
