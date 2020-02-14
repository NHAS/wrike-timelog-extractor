package models

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
