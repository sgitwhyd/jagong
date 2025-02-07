package router

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/app/repository"
	"github.com/sgitwhyd/jagong/pkg/jwt_token"
	"github.com/sgitwhyd/jagong/pkg/response"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	accessToken := ctx.Get("Authorization")
	if accessToken == "" {
		log.Printf("middleware: token is empty")
		err := errors.New("token is empty").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	_, err := repository.GetAuthSessionByToken(ctx.Context(), accessToken)
	if err != nil {
		log.Printf("get auth session err:%v", err)
		err := errors.New("token is invalid").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	claimsToken, err := jwt_token.ValidateToken(ctx.Context(), accessToken)
	if err != nil {
		log.Printf("middleware: token validation error: %v", err)
		err := errors.New("unauthorized").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	//check if jwt_token expired time is still active or not
	if time.Now().Unix() > claimsToken.ExpiresAt.Unix() {
		err := errors.New("token is expired").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	//set claim jwt_token data into context
	ctx.Locals("username", claimsToken.Username)
	ctx.Locals("full_name", claimsToken.FullName)

	return ctx.Next()
}

func RefreshTokenMiddleware(ctx *fiber.Ctx) error {
	refreshToken := ctx.Get("Refresh-Token")
	if refreshToken == "" {
		err := errors.New("refresh token is empty").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	claimsToken, err := jwt_token.ValidateToken(ctx.Context(), refreshToken)
	if err != nil {
		log.Printf("middleware: refresh token validation error: %v", err)
		err := errors.New("unauthorized").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	//check if jwt_token expired time is still active or not
	if time.Now().Unix() > claimsToken.ExpiresAt.Unix() {
		err := errors.New("refresh token is expired").Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnauthorized, &err)
	}

	//set claim jwt_token data into context
	ctx.Locals("username", claimsToken.Username)
	ctx.Locals("full_name", claimsToken.FullName)

	return ctx.Next()
}
