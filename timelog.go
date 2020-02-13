package main

import (
	"encoding/json"
)

type Timelog struct {
	User        string `json:"-"` // Not set in unmarshalling, but found from contacts map
	UserId      string
	TaskId      string
	FirstName   string
	LastName    string
	TrackedDate string
	Hours       float64
}

type Timelogs struct {
	Data []Timelog
}

func getTimelogMap(apiKey string, contacts map[string]Contact) (result map[string][]Timelog, err error) {
	textContent := getDataForURL(host+"/timelogs", apiKey)

	var timelogs Timelogs
	err = json.Unmarshal(textContent, &timelogs)
	if err != nil {
		return result, err
	}

	result = make(map[string][]Timelog)
	for _, k := range timelogs.Data {

		k.User = contacts[k.UserId].FirstName + " " + contacts[k.UserId].LastName
		result[k.TaskId] = append(result[k.TaskId], k)
	}

	return result, nil
}

type taskTimeLog struct {
	fields  map[string]string
	timelog Timelog
}

type timeLogCustomFields struct {
	Id    string
	Value string
}

type collectiveTimeLog struct {
	Id           string
	CustomFields []timeLogCustomFields
}

type collectiveTimelogs struct {
	Data []collectiveTimeLog
}

func getTaskTimelogs(apiKey string, customFields map[string]string, timelogs map[string][]Timelog) (result []taskTimeLog, err error) {
	textContent := getDataForURL(host+"/tasks?subTasks=true&fields=['customFields']", apiKey)

	var cTimelogs collectiveTimelogs
	err = json.Unmarshal(textContent, &cTimelogs)
	if err != nil {
		return result, err
	}

	for _, k := range cTimelogs.Data {

		fields := make(map[string]string)
		for _, field := range k.CustomFields {
			if field.Value != "" {
				fields[field.Id] = field.Value
			}
		}

		if len(fields) > 0 {
			for key := range customFields {
				if _, exists := fields[key]; !exists {
					fields[key] = ""
				}
			}

			taskLogs := timelogs[k.Id]
			for _, log := range taskLogs {
				var task = taskTimeLog{}
				task.fields = fields
				task.timelog = log
				result = append(result, task)
			}
		}
	}

	return result, nil
}
