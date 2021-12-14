package google

import "context"

type GoogleRepository interface {
	CountGoogleUserByEmail(ctx context.Context, email string) (total int, err error)
}
