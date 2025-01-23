package repository

import (
	"context"
	"github.com/sgitwhyd/jagong/app/models"
	"github.com/sgitwhyd/jagong/pkg/database"
	"time"
)

func CreateUser(ctx context.Context, model *models.User) error {
	return database.DB.Create(model).Error
}

func FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	model := &models.User{}
	err := database.DB.Where("username = ?", username).First(model).Error
	return model, err
}

func UpdateAuthSessionByUserId(ctx context.Context, userID uint, newToken string, newTokenExpiry time.Time) error {
	return database.DB.Exec("UPDATE user_sessions SET token = ?, token_expired = ? WHERE user_id = ?", newToken, newTokenExpiry, userID).Error
}

func GetAuthSessionByToken(ctx context.Context, token string) (*models.UserSession, error) {
	model := &models.UserSession{}
	err := database.DB.Where("token = ?", token).First(model).Error
	return model, err
}

func CreateAuthSession(ctx context.Context, model models.UserSession) error {
	return database.DB.Create(&model).Error
}

func DeleteAuthSessionByToken(ctx context.Context, token string) error {
	return database.DB.Where("token = ?", token).Delete(&models.UserSession{}).Error
}
