package response

import (
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Error   *string     `json:"error"`
}

func SendSuccessResponse(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(Response{
		Message: "success",
		Code:    fiber.StatusOK,
		Data:    data,
		Error:   nil,
	})
}

func SendErrorResponse(ctx *fiber.Ctx, code int, err *string) error {
	return ctx.Status(code).JSON(&Response{
		Message: "Failed",
		Code:    code,
		Data:    nil,
		Error:   err,
	})
}
