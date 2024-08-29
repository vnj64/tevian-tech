package upload_task_image

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"tevian/domain"
)

type Request struct {
	Id        string `json:"id"`
	ImageName string `json:"imageName"`
	Body      []byte `json:"body"`
}

func Run(c domain.Context, r Request) error {
	_, err := c.Connection().Task().WhereId(r.Id)
	if err != nil {
		return fmt.Errorf("task with id [%s] does not exist due [%v]", r.Id, err)
	}

	r.ImageName = uuid.New().String() + ".jpg"
	file, err := os.Create(fmt.Sprintf("images/%s", r.ImageName))
	if err != nil {
		return fmt.Errorf("failed to create file due [%v]", err)
	}
	defer file.Close()

	reader := bytes.NewReader(r.Body)

	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("failed to write to file due [%v]", err)
	}

	updates := make(map[string]interface{})

	updates["image_name"] = r.ImageName
	updates["image_address"] = fmt.Sprintf("images/%s", r.ImageName)

	if len(updates) == 0 {
		return fmt.Errorf("nothing to update")
	}

	if err := c.Connection().Task().Update(r.Id, updates); err != nil {
		return fmt.Errorf("failed to update task due [%v]", err)
	}

	return nil
}
