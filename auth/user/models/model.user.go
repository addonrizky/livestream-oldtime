package models

import "time"

type UserReq struct {
	ID              int64      `json:"id" gorm:"column:id"`
	UserName        string     `json:"username"  gorm:"column:username"`
	Password        string     `json:"password"  gorm:"column:password"`
	Name            string     `json:"name"  gorm:"column:name"`
	LoginType       string     `json:"login_type"  gorm:"column:login_type"`
	Email           string     `json:"email"  gorm:"column:email"`
	ConfirmPassword string     `json:"confirm_password"`
	RoleID          int64      `json:"role_id"  gorm:"column:role_id"`
	IsActive        bool       `json:"is_active"  gorm:"column:is_active"`
	RegisterDate    time.Time  `json:"register_date"  gorm:"column:register_date"`
	UpdatedDate     *time.Time `json:"updated_date"  gorm:"column:updated_date"`
}
type User struct {
	ID           int64      `json:"id" gorm:"column:id"`
	UserName     string     `json:"username"  gorm:"column:username"`
	Password     string     `json:"password"  gorm:"column:password"`
	Name         string     `json:"name"  gorm:"column:name"`
	LoginType    string     `json:"login_type"  gorm:"column:login_type"`
	Email        string     `json:"email"  gorm:"column:email"`
	RoleID       int64      `json:"role_id"  gorm:"column:role_id"`
	IsActive     bool       `json:"is_active"  gorm:"column:is_active"`
	RegisterDate time.Time  `json:"register_date"  gorm:"column:register_date"`
	UpdatedDate  *time.Time `json:"updated_date"  gorm:"column:updated_date"`
}

type ResetPasswordReq struct {
	Email           string `json:"email" `
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  bool   `json:"result"`
}

type EmailResetPasswordReq struct {
	Email string `json:"email" `
}

type EmailResetPasswordResponse struct {
	Code              int    `json:"code"`
	Message           string `json:"message"`
	ResetPasswordLink string `json:"reset_password_link"`
}
