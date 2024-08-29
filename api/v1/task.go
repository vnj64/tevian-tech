package v1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"io/ioutil"
	"tevian/domain"
	"tevian/domain/cases/create_task"
	"tevian/domain/cases/delete_task"
	"tevian/domain/cases/start_task"
	"tevian/domain/cases/upload_task_image"
)

func CreateTaskHandler(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	var request create_task.Request

	uid := uuid.New().String()
	request.Task.Id = uid

	if err := ctx.BodyParser(&request); err != nil {
		return BadRequest(err)
	}

	id, err := create_task.Run(c, request)
	if err != nil {
		return InternalServerError(err)
	}

	return OK(id)
}

func UploadTaskImageHandler(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	form, err := ctx.MultipartForm()
	if err != nil {
		return BadRequest(err)
	}

	fileHeader := form.File["body"][0]
	file, err := fileHeader.Open()
	if err != nil {
		return InternalServerError(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return InternalServerError(err)
	}

	request := upload_task_image.Request{
		Id:        ctx.Params("id"),
		ImageName: uuid.New().String(),
		Body:      fileBytes,
	}

	if err := upload_task_image.Run(c, request); err != nil {
		return InternalServerError(err)
	}

	return OK(nil)
}

func DeleteTaskHandler(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	err := delete_task.Run(c, delete_task.Request{
		Id: ctx.Params("id"),
	})

	// TODO: надо будет что-то другое придумать
	if errors.Is(err, errors.New("cant delete task, which still processing")) {
		return Forbidden(err)
	}

	if err = delete_task.Run(c, delete_task.Request{
		Id: ctx.Params("id"),
	}); err != nil {
		return InternalServerError(err)
	}

	return OK(nil)
}

func StartTaskHandler(c domain.Context, ctx *fiber.Ctx) *RawResponse {
	result, err := start_task.Run(c, start_task.Request{
		Id: ctx.Params("id"),
	})
	if err != nil {
		return InternalServerError(err)
	}

	return OK(result)
}
