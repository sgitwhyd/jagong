package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/app/repository"
	"github.com/sgitwhyd/jagong/pkg/response"
	"go.elastic.co/apm/v2"
)

func GetHistory(ctx *fiber.Ctx) error {
	context := ctx.Context()
	span , spanCtx := apm.StartSpan(context, "GetHistory", "controller")
	defer span.End()
	
	resp, err := repository.FindAllMessage(spanCtx)
	if err != nil {
		log.Printf("get history message error %v", err.Error())
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	return response.SendSuccessResponse(ctx, resp)
}
