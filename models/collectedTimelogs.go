package models

type CollectiveTimeLog struct {
	Id           string
	SubTaskIds   []string
	ParentIds    []string
	CustomFields []CustomField
}

type CollectiveTimelogs struct {
	Data []CollectiveTimeLog
}
