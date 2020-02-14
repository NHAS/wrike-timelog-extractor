package models

type CustomField struct {
	Id    string
	Value string
	Title string
}

type CustomFields struct {
	Data []CustomField
}
