package v1

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"tevian/domain"
)

type RawResponse struct {
	error      error
	status     int
	additional interface{}
	payload    interface{}
}

func (r *RawResponse) Body() *ResponseBody {
	return &ResponseBody{
		Response: Response{
			Status: r.status,
		},
		Additional: r.additional,
		Payload:    r.payload,
	}
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

type ResponseBody struct {
	Response   `json:"response"`
	Additional interface{} `json:"additional,omitempty"`
	Payload    interface{} `json:"payload,omitempty"`
}

func WrapHandler(handler func(c domain.Context, ctx *fiber.Ctx) *RawResponse) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newContext, ok := ctx.Locals("context").(domain.Context)
		if !ok {
			return nil
		}

		response := handler(newContext, ctx)
		body := response.Body()

		status := body.Status

		if err := response.Error(); err != nil {

			body.Message = response.Error().Error()
		}
		return ctx.Status(status).JSON(body)
	}
}

func (r *RawResponse) Error() error {
	return r.error
}

func BadRequest(err error) *RawResponse {
	return &RawResponse{
		status: http.StatusBadRequest,
		error:  err,
	}
}

func Forbidden(err error) *RawResponse {
	return &RawResponse{
		status: http.StatusForbidden,
		error:  err,
	}
}

func (r *RawResponse) WithExtraCode(code int) *RawResponse {
	r.status = code
	return r
}

func InternalServerError(err error) *RawResponse {
	return &RawResponse{
		status: http.StatusInternalServerError,
		error:  err,
	}
}

func (r *RawResponse) WithPayload(payload any) *RawResponse {
	r.payload = payload
	return r
}

func OK(payload any) *RawResponse {
	out := &RawResponse{
		status: http.StatusOK,
	}

	if payload != nil {
		out.WithPayload(payload)
	}

	return out
}
