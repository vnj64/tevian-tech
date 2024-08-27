package services

import "tevian/domain/models"

type DetectFaceResponse struct {
	TaskId string                     `json:"taskId"`
	Status *models.TaskStatus         `json:"status"`
	Images []ImagesAdditional         `json:"images"`
	Stats  map[string]StatsAdditional `json:"stats"`
}

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

type ImagesAdditional struct {
	ImageName string            `json:"imageName"`
	Faces     []FacesAdditional `json:"faces"`
}

type FacesAdditional struct {
	FaceBbox   string `json:"faceBbox"`
	FaceGender Gender `json:"gender"`
	Age        uint32 `json:"age"`
}

type StatsAdditional struct {
	AllFacesQuantity int     `json:"allFacesQuantity"`
	MaleQuantity     int     `json:"maleQuantity"`
	FemaleQuantity   int     `json:"femaleQuantity"`
	AverageMaleAge   float32 `json:"averageMaleAge"`
	AverageFemaleAge float32 `json:"averageFemaleAge"`
}

type FaceCloud interface {
	Detect(jwtToken string) (DetectFaceResponse, error)
}
