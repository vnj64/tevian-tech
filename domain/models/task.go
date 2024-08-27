package models

const (
	StatusForming    TaskStatus = "FORMING"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

type TaskStatus string

type Task struct {
	Id     string     `json:"id"`
	Status TaskStatus `json:"status"`
}
