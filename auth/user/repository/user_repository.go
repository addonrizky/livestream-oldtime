package repository

import (
	"context"
	"time"

	user "github.com/asumsi/livestream/auth/user"
	"github.com/asumsi/livestream/auth/user/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	Conn *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.UserRepository {
	return &UserRepository{
		Conn: db,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, req *models.UserReq) error {
	user := models.User{UserName: req.UserName, Password: req.Password, Name: req.Name, LoginType: req.LoginType, Email: req.Email, RoleID: req.RoleID, IsActive: req.IsActive, RegisterDate: time.Now()}

	sql := r.Conn.WithContext(ctx).Table("users").Create(&user)
	if sql.Error != nil {
		return sql.Error
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (res models.User, err error) {
	sql := r.Conn.WithContext(ctx).Table("users").Select("id, username, name, password, login_type, is_active").Where("email = ?", email).Scan(&res)
	if sql.Error != nil {
		return res, sql.Error
	}
	return res, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (res models.User, err error) {
	sql := r.Conn.WithContext(ctx).Table("users").Select("id, username, name, password, login_type, is_active").Where("username = ?", username).Scan(&res)
	if sql.Error != nil {
		return res, sql.Error
	}
	return res, nil
}

func (r *UserRepository) UpdatePasswordByEmail(ctx context.Context, req *models.ResetPasswordReq) error {
	now := time.Now()
	user := &models.User{Email: req.Email, Password: req.Password, UpdatedDate: &now}

	sql := r.Conn.WithContext(ctx).Table("users").
		Where("email = ?", user.Email).
		Where("is_active = ?", true).
		Where("login_type = ?", "APP").
		Updates(&user)
	if sql.Error != nil {
		return sql.Error
	}

	return nil
}
