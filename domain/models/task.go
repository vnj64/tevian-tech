package models

const (
	StatusForming    TaskStatus = "FORMING"
	StatusProcessing TaskStatus = "PROCESSING"
	StatusCompleted  TaskStatus = "COMPLETED"
	StatusFailed     TaskStatus = "FAILED"
)

type TaskStatus string

type Task struct {
	Id               string     `json:"id"`
	Status           TaskStatus `json:"status"`
	AllFacesQuantity *int
	MaleQuantity     *int
	FemaleQuantity   *int
	AverageMaleAge   *float64
	AverageFemaleAge *float64
}

type BoundingBox struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Faces struct {
	BoundingBox BoundingBox `json:"boundingBox"`
	Gender      string      `json:"gender"`
	Age         float64     `json:"age"`
}

type ImageData struct {
	Name  string  `json:"name"`
	Faces []Faces `json:"faces"`
}

type Statistics struct {
	TotalFaces       int     `json:"totalFaces"`
	TotalMales       int     `json:"totalMales"`
	TotalFemales     int     `json:"totalFemales"`
	AverageMaleAge   float64 `json:"averageMaleAge"`
	AverageFemaleAge float64 `json:"averageFemaleAge"`
}

type DetectResult struct {
	TaskId     string      `json:"taskId"`
	Status     TaskStatus  `json:"status"`
	ImageData  []ImageData `json:"imageData"`
	Statistics Statistics  `json:"statistics"`
}
