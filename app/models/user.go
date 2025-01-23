package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key"`
	Username  string `gorm:"unique;type:varchar(20)" json:"username" validate:"required,min=6,max=20"`
	FullName  string `gorm:"type:varchar(255)" json:"full_name" validate:"required,min=6,max=255"`
	Password  string `gorm:"type:varchar(255);" json:"password"  validate:"required,min=6"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserSession struct {
	ID                  uint `gorm:"primary_key"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserID              uint      `gorm:"index" json:"user_id" validate:"required,min=1,required"`
	Token               string    `gorm:"type:varchar(255)" json:"jwt_token" validate:"required`
	RefreshToken        string    `gorm:"type:varchar(255)" json:"refresh_token" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l UserLoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserLoginResponse struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Token        string `json:"jwt_token"`
	RefreshToken string `json:"refresh_token"`
}
