package models

type Image struct {
	Id           string `json:"id"`
	TaskId       string `json:"task"`
	ImageName    string `json:"imageName"`
	ImageAddress string `json:"imageAddress"`
}
