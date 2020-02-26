package main

import (
	"flag"
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

	apiKeyArg := flag.String("apiKey", "", "A Wrike API key, created through the UI (required); can also be set via the env var WRIKEKEY")
	startDateArg := flag.String("start-date", "", "an explicit date to filter from (optional); if end date is not set will be until now.")
	endDateArg := flag.String("end-date", "", "an explicit date to filter to (optional); requires start-date be set")
	recentArg := flag.String("for-last", "month", "either day, week, month (default) or quarter, and is the same as setting a start date at now - duration; overriden by specific dates")

	flag.Parse()

	validateArgs(apiKey, apiKeyArg, startDateArg, endDateArg, recentArg)

	if *apiKeyArg != "" {
		apiKey = *apiKeyArg
	}

	tasks, customFields := gatherData(apiKey)
	csv := asCsv(tasks, customFields)

	fmt.Print(csv)
}

func validateArgs(apiKey string, apiKeyArg *string, startDateArg *string, endDateArg *string, recentArg *string) {
	if apiKey == "" && *apiKeyArg == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *startDateArg == "" && *endDateArg != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if (*startDateArg != "" || *endDateArg != "") && *recentArg != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
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
