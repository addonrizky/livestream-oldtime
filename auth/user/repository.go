package user

import (
	"context"

	"github.com/asumsi/livestream/auth/user/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req *models.UserReq) error
	GetUserByEmail(ctx context.Context, email string) (res models.User, err error)
	GetUserByUsername(ctx context.Context, username string) (res models.User, err error)
	UpdatePasswordByEmail(ctx context.Context, req *models.ResetPasswordReq) error
}
