package domain

import (
	"tevian/domain/repositories"
)

type Connection interface {
	Task() repositories.Task
	Face() repositories.Face
	Image() repositories.Image
}
