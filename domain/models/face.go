package models

type Face struct {
	Id      string  `json:"id"`
	ImageId string  `json:"imageId"`
	Bbox    string  `json:"bbox"`
	Gender  string  `json:"gender"`
	Age     float64 `json:"age"`
}
