package google

import "context"

type GoogleUsecase interface {
	CheckIsUserAlreadyRegistered(ctx context.Context, email string) (bool, error)
}
