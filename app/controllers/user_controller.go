package controllers

import (
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/app/repository"
	"github.com/sgitwhyd/jagong/pkg/jwt_token"
	"github.com/sgitwhyd/jagong/pkg/response"
	"go.elastic.co/apm/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *fiber.Ctx) error {
	context := ctx.Context()

	span , spanCtx := apm.StartSpan(context, "Register", "controller")
	defer span.End()

	user := new(models.User)

	err := ctx.BodyParser(&user)
	if err != nil {
		log.Printf("parse body err:%v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnprocessableEntity, &err)
	}

	err = user.Validate()
	if err != nil {
		log.Printf("validate err:%v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusBadRequest, &err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("bcrypt hash err:%v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}
	user.Password = string(hashedPassword)

	err = repository.CreateUser(spanCtx, user)
	if err != nil {
		log.Printf("create user error: %v", err.Error())
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	return response.SendSuccessResponse(ctx, nil)
}

func Login(ctx *fiber.Ctx) error {
	context := ctx.Context()
	span , spanCtx := apm.StartSpan(context, "Login", "controller")
	defer span.End()

	loginRequest := new(models.UserLoginRequest)

	err := ctx.BodyParser(&loginRequest)
	if err != nil {
		log.Printf("parse body err:%v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusUnprocessableEntity, &err)
	}

	err = loginRequest.Validate()
	if err != nil {
		log.Printf("validate err:%v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusBadRequest, &err)
	}

	user, err := repository.FindUserByUsername(spanCtx, loginRequest.Username)
	if err != nil {
		log.Printf("find user error: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := "Combine username and password does not exist"
			return response.SendErrorResponse(ctx, fiber.StatusNotFound, &err)
		}
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		log.Printf("compare hash err:%v", err)
		err := "Combine username and password does not exist"
		return response.SendErrorResponse(ctx, fiber.StatusNotFound, &err)
	}

	generatedToken, err := jwt_token.GenerateToken(spanCtx, user.Username, user.FullName, "token")
	if err != nil {
		log.Printf("generate token error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	generatedRefreshToken, err := jwt_token.GenerateToken(spanCtx, user.Username, user.FullName, "refresh_token")
	if err != nil {
		log.Printf("generate refresh token error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	authSessionData := models.UserSession{
		UserID:              user.ID,
		Token:               generatedToken,
		TokenExpired:        time.Now().Add(jwt_token.MapTypeToken["token"]),
		RefreshToken:        generatedRefreshToken,
		RefreshTokenExpired: time.Now().Add(jwt_token.MapTypeToken["refresh_token"]),
	}

	err = repository.CreateAuthSession(spanCtx, authSessionData)
	if err != nil {
		log.Printf("create auth session error: %v", err)
		err := errors.New("login failed").Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	loginResponse := models.UserLoginResponse{
		Username:     user.Username,
		FullName:     user.FullName,
		Token:        generatedToken,
		RefreshToken: generatedRefreshToken,
	}

	return response.SendSuccessResponse(ctx, loginResponse)
}

func Logout(ctx *fiber.Ctx) error {
	context := ctx.Context()
	span , spanCtx := apm.StartSpan(context, "Logout", "controller")
	defer span.End()

	userToken := ctx.Get("Authorization")

	err := repository.DeleteAuthSessionByToken(spanCtx, userToken)
	if err != nil {
		log.Printf("delete auth session error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	return response.SendSuccessResponse(ctx, nil)
}

func RefreshToken(ctx *fiber.Ctx) error {
	context := ctx.Context()
	span , spanCtx := apm.StartSpan(context, "RefreshToken", "controller")
	defer span.End()
	username := ctx.Locals("username").(string)

	user, err := repository.FindUserByUsername(context, username)
	if err != nil {
		log.Printf("find user error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	generatedToken, err := jwt_token.GenerateToken(spanCtx, username, username, "token")
	if err != nil {
		log.Printf("generate token error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	err = repository.UpdateAuthSessionByUserId(spanCtx, user.ID, generatedToken, time.Now().Add(jwt_token.MapTypeToken["token"]))
	if err != nil {
		log.Printf("update auth session error: %v", err)
		err := err.Error()
		return response.SendErrorResponse(ctx, fiber.StatusInternalServerError, &err)
	}

	refreshTokenResponse := fiber.Map{
		"token": generatedToken,
	}

	return response.SendSuccessResponse(ctx, refreshTokenResponse)
}
