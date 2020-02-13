package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const host = "https://www.wrike.com/api/v4"

func getDataForURL(url string, apiKey string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func main() {
	log.SetFlags(0) // Disable log timestamp

	apiKey := os.Getenv("WRIKEKEY")
	if len(os.Args) == 2 {
		apiKey = os.Args[1]
	}

	if apiKey == "" {
		fmt.Println("please provide an API key (permananent access token) as an argument or as the env var WRIKEKEY")
		os.Exit(1)
	}

	customFields, err := getCustomFieldsMap(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	contacts, err := getContactsMap(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	timeLogs, err := getTimelogMap(apiKey, contacts)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := getTaskTimelogs(apiKey, customFields, timeLogs)
	if err != nil {
		log.Fatal(err)
	}

	csv := ""
	for _, task := range tasks {
		for key := range customFields {
			csv += task.fields[key] + ","
		}
		csv += fmt.Sprintf("%s,%s,%f\n", task.timelog.User, task.timelog.TrackedDate, task.timelog.Hours)
	}

	fmt.Print(csv)
}
