package usecase

import (
	"context"
	"errors"
	"time"

	user "github.com/asumsi/livestream/auth/user"
	"github.com/asumsi/livestream/auth/user/models"
	"github.com/asumsi/livestream/modelgeneral"
	"github.com/asumsi/livestream/utility"
	"github.com/dgrijalva/jwt-go"
)

type userUsecase struct {
	userRepo       user.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(userRepo user.UserRepository, timeout time.Duration) user.UserUsecase {
	return &userUsecase{
		userRepo:       userRepo,
		contextTimeout: timeout,
	}
}

func (a *userUsecase) CreateUser(ctx context.Context, req *models.UserReq) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	if req.Password != req.ConfirmPassword {
		return errors.New("password is not same")
	}

	hashPass, err := utility.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashPass

	err = a.userRepo.CreateUser(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func (a *userUsecase) Authenticate(ctx context.Context, req modelgeneral.LoginRequest, loginType string) (resp models.AuthResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	var user models.User

	user.Email = req.Email
	user.UserName = req.Email
	user.RoleID = int64(1)

	if loginType == "APP" {
		if utility.IsEmailValid(req.Email) {
			user, err = a.userRepo.GetUserByEmail(ctx, req.Email)
		} else {
			user, err = a.userRepo.GetUserByUsername(ctx, req.Username)
		}
		if user.ID == 0 || err != nil {
			return resp, errors.New("Username/Email tidak terdaftar")
		}
		if !user.IsActive {
			return resp, errors.New("Username/Email tidak aktif")
		}

		if !utility.CheckPasswordHash(req.Password, user.Password) {
			return resp, errors.New("password salah")
		}
	}

	resultsAll, err := a.GenerateToken(ctx, user.UserName, user.Email, user.RoleID)
	return resultsAll, err
}

func (a *userUsecase) GenerateToken(ctx context.Context, name, email string, roleID int64) (resp models.AuthResponse, err error) {
	appName := utility.GetConfigString(`app_name`)
	jwtMethod := utility.GetConfigString(`jwt_method`)
	jwtSecret := utility.GetConfigString(`jwt_secret`)
	jwtLifespan := utility.GetConfigDuration(`jwt_lifespan`)
	jwtLifespanRefresh := utility.GetConfigDuration(`jwt_refresh_lifespan`)
	resp.Token, err = utility.GenerateJwtToken(modelgeneral.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(jwtLifespan).Unix(),
		},
		Name:   name,
		Email:  email,
		RoleID: roleID,
	}, jwtMethod, jwtSecret)
	resp.RefreshToken, err = utility.GenerateJwtToken(modelgeneral.JwtClaims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    appName,
			ExpiresAt: time.Now().Add(jwtLifespanRefresh).Unix(),
		},
		Name:  name,
		Email: email,
	}, jwtMethod, jwtSecret)
	return resp, nil
}

func (a *userUsecase) UpdatePassword(ctx context.Context, req *models.ResetPasswordReq) error {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	if req.Password != req.ConfirmPassword {
		return errors.New("password is not same")
	}

	hashPass, err := utility.HashPassword(req.Password)
	if err != nil {
		return err
	}

	req.Password = hashPass

	err = a.userRepo.UpdatePasswordByEmail(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
