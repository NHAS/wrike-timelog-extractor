package models

type Contact struct {
	Id        string
	FirstName string
	LastName  string
}

type Contacts struct {
	Data []Contact
}
