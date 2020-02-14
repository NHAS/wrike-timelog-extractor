package main

import (
	"encoding/json"
)

type taskTimeLog struct {
	fields  map[string]string
	timelog Timelog
}

func findCustomFields(task collectiveTimeLog, parentMap map[string]collectiveTimeLog) map[string]string {
	if len(task.CustomFields) == 0 {
		if _, exists := parentMap[task.Id]; !exists {
			return make(map[string]string)
		}

		parent := parentMap[task.Id]
		return findCustomFields(parent, parentMap)
	}

	fields := make(map[string]string)
	for _, field := range task.CustomFields {
		if field.Value != "" {
			fields[field.Id] = field.Value
		}
	}

	return fields
}

func getTaskTimelogs(apiKey string, customFields map[string]CustomField, timelogs map[string][]Timelog, folders map[string]collectiveTimeLog) (result []taskTimeLog, err error) {
	textContent := getDataForURL(host+"/tasks?subTasks=true&fields=['customFields','subTaskIds','parentIds']", apiKey)

	var cTimelogs collectiveTimelogs
	err = json.Unmarshal(textContent, &cTimelogs)
	if err != nil {
		return result, err
	}

	// used to collect closest parent custom fields if necessary
	parentMap := make(map[string]collectiveTimeLog)
	for _, k := range cTimelogs.Data {
		for _, j := range k.SubTaskIds {
			parentMap[j] = k
		}
		for _, j := range k.ParentIds {
			if f, exists := folders[j]; exists {
				parentMap[k.Id] = f
			}
		}
	}

	for _, k := range cTimelogs.Data {

		fields := findCustomFields(k, parentMap)

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
