package models

type Image struct {
	Id     string `json:"id"`
	TaskId string `json:"task"`
	Name   string `json:"name"`
}
