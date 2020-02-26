package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const host = "https://www.wrike.com/api/v4"

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

	tasks, customFields := gatherData(apiKey)
	csv := asCsv(tasks, customFields)

	fmt.Print(csv)
}

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

func gatherData(apiKey string) ([]taskTimeLog, []string) {
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

	folders, err := getFoldersAsTasks(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err := getTaskTimelogs(apiKey, customFields, timeLogs, folders)
	if err != nil {
		log.Fatal(err)
	}

	customFieldKeys := make([]string, 0)
	for key := range customFields {
		customFieldKeys = append(customFieldKeys, key)
	}

	return tasks, customFieldKeys
}

func asCsv(tasks []taskTimeLog, customFields []string) string {
	csv := ""

	for _, task := range tasks {
		for _, key := range customFields {
			csv += task.fields[key] + ","
		}
		csv += fmt.Sprintf("%s,%s,%f\n", task.timelog.User, task.timelog.TrackedDate, task.timelog.Hours)
	}

	return csv
}
