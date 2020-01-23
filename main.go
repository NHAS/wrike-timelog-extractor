package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getDataForURL(url string, apiKey string) interface{} {
	var bearer = "Bearer " + apiKey
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)

	var f interface{}
	json.Unmarshal([]byte(body), &f)

	return f
}

func main() {
	fmt.Println("wrike-extractor running!")

	if len(os.Args) != 2 {
		fmt.Println("please provide an API key (permananent access token) as an argument")
		os.Exit(1)
	}

	var apiKey = os.Args[1]
	var host = "https://www.wrike.com/api/v4"

	var customFields = getDataForURL(host+"/customfields", apiKey)
	// var tasks = getDataForURL(host+"/tasks?fields=['customFields']", apiKey)
	// var timelogs = getDataForURL(host+"/timelogs", apiKey)
	// var contacts = getDataForURL(host+"/contacts", apiKey)

	fmt.Println(customFields)

	// TODO: compose CSV and print out
}
