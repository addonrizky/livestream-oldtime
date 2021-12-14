package user

import (
	"context"

	"github.com/asumsi/livestream/auth/user/models"
	"github.com/asumsi/livestream/modelgeneral"
)

type UserUsecase interface {
	CreateUser(ctx context.Context, req *models.UserReq) error
	Authenticate(ctx context.Context, req modelgeneral.LoginRequest, loginType string) (resp models.AuthResponse, err error)
	GenerateToken(ctx context.Context, name, email string, roleID int64) (resp models.AuthResponse, err error)
	UpdatePassword(ctx context.Context, req *models.ResetPasswordReq) error
}
