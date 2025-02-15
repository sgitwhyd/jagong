package repository

import (
	"context"
	"time"

	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/pkg/database"
	"go.elastic.co/apm/v2"
)

func CreateUser(ctx context.Context, model *models.User) error {
	span, _ := apm.StartSpan(ctx, "CreateUser", "repository")
	defer span.End()

	return database.DB.Create(model).Error
}

func FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	span, _ := apm.StartSpan(ctx, "FindUserByUsername", "repository")
	defer span.End()

	model := &models.User{}
	err := database.DB.Where("username = ?", username).First(model).Error
	return model, err
}

func UpdateAuthSessionByUserId(ctx context.Context, userID uint, newToken string, newTokenExpiry time.Time) error {
	span, _ := apm.StartSpan(ctx, "UpdateAuthSessionByUserId", "repository")
	defer span.End()
	return database.DB.Exec("UPDATE user_sessions SET token = ?, token_expired = ? WHERE user_id = ?", newToken, newTokenExpiry, userID).Error
}

func GetAuthSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	span, _ := apm.StartSpan(ctx, "GetAuthSessionByToken", "repository")
	defer span.End()

	model := &models.UserSession{}
	err := database.DB.Where("token = ?", token).First(model).Error
	return model, err
}

func CreateAuthSession(ctx context.Context, model models.UserSession) error {
	span, _ := apm.StartSpan(ctx, "CreateAuthSession", "repository")
	defer span.End()
	return database.DB.Create(&model).Error
}

func DeleteAuthSessionByToken(ctx context.Context, token string) error {
	span, _ := apm.StartSpan(ctx, "DeleteAuthSessionByToken", "repository")
	defer span.End()
	return database.DB.Where("token = ?", token).Delete(&models.UserSession{}).Error
}
