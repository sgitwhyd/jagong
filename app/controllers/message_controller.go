package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/app/repository"
	"github.com/sgitwhyd/jagong/pkg/response"
)

func GetMessages(ctx *fiber.Ctx) error {
	context := ctx.Context()
	resp, err := repository.FindAllMessage(context)
	if err != nil {
		fmt.Printf("get history message error %v", err.Error())
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	return response.SendSuccessResponse(ctx, resp)
}
