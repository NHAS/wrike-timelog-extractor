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

	// var tasks = getDataForURL(host+"/tasks?fields=['customFields']", apiKey)
	// var timelogs = getDataForURL(host+"/timelogs", apiKey)

	fmt.Println(customFields)
	fmt.Println(contacts)

	// TODO: compose CSV and print out
}
