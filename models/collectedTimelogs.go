package models

type collectiveTimeLog struct {
	Id           string
	SubTaskIds   []string
	ParentIds    []string
	CustomFields []CustomField
}

type collectiveTimelogs struct {
	Data []collectiveTimeLog
}
