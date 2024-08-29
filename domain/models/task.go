package models

const (
	StatusForming    TaskStatus = "FORMING"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
)

type TaskStatus string
type Gender string

type Task struct {
	Id           string     `json:"id"`
	Status       TaskStatus `json:"status"`
	ImageAddress string     `json:"imageAddress"`
	ImageName    string     `json:"imageName"`
}

type ResultData struct {
	TaskId string     `json:"taskId"`
	Status TaskStatus `json:"status"`
	Images struct {
		Name  string `json:"imageName"`
		Faces []struct {
			BBox struct {
				Height int `json:"height"`
				Width  int `json:"width"`
				X      int `json:"x"`
				Y      int `json:"y"`
			} `json:"bbox"`
			FaceGender Gender `json:"faceGender"`
			Age        uint32 `json:"age"`
		} `json:"faces"`
	} `json:"images"`
	Stats struct {
		AllFacesQuantity int     `json:"allFacesQuantity"`
		AverageGender    int     `json:"averageGender"`
		AverageMaleAge   float32 `json:"averageMaleAge"`
		AverageFemaleAge float32 `json:"averageFemaleAge"`
	} `json:"stats"`
}
