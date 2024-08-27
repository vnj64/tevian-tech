package domain

import (
	"tevian/domain/repositories"
)

type Connection interface {
	Task() repositories.Task
}
