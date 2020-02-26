package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const host = "https://www.wrike.com/api/v4"

func main() {
	log.SetFlags(0) // Disable log timestamp

	apiKey, start, end := getConfig()
	tasks, customFields := gatherData(apiKey, start, end)
	csv := asCsv(tasks, customFields)

	fmt.Print(csv)
}

func getConfig() (apiKey string, start time.Time, end time.Time) {

	apiKeyEnv := os.Getenv("WRIKEKEY")
	apiKeyArg := flag.String("api-key", "", "A Wrike API key, created through the UI (required); can also be set via the env var WRIKEKEY")
	startDateArg := flag.String("start-date", "", "an explicit date to filter from in the format dd/MM/yyyy (optional); if end-date is not set will be until now.; if start-date is not set, assumes last month")
	endDateArg := flag.String("end-date", "", "an explicit date to filter to in the format dd/MM/yyyy (optional); requires start-date be set")

	flag.Parse()

	if apiKeyEnv == "" && *apiKeyArg == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	apiKey = apiKeyEnv
	if *apiKeyArg != "" {
		apiKey = *apiKeyArg
	}

	if *startDateArg == "" && *endDateArg != "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	now := time.Now()

	if *startDateArg == "" {
		return apiKey, now.AddDate(0, -1, 0), now
	}

	start, err := time.Parse("02/01/2006", *startDateArg)
	if err != nil || start.After(now) {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *endDateArg != "" {
		end, err = time.Parse("02/01/2006", *endDateArg)
		if err != nil || end.Before(start) {
			flag.PrintDefaults()
			os.Exit(1)
		}
	} else {
		end = now
	}

	return apiKey, start, end
}

func gatherData(apiKey string, start time.Time, end time.Time) (tasks []taskTimeLog, customFieldKeys []string) {
	customFields, err := getCustomFieldsMap(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	contacts, err := getContactsMap(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	timeLogs, err := getTimelogMap(apiKey, contacts, start, end)
	if err != nil {
		log.Fatal(err)
	}

	folders, err := getFoldersAsTasks(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	tasks, err = getTaskTimelogs(apiKey, customFields, timeLogs, folders)
	if err != nil {
		log.Fatal(err)
	}

	customFieldKeys = make([]string, 0)
	for key := range customFields {
		customFieldKeys = append(customFieldKeys, key)
	}

	return tasks, customFieldKeys
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
