package repository

import (
	"context"

	google "github.com/asumsi/livestream/auth/google"
	"gorm.io/gorm"
)

type GoogleRepository struct {
	Conn *gorm.DB
}

func NewGoogleRepository(db *gorm.DB) google.GoogleRepository {
	return &GoogleRepository{
		Conn: db,
	}
}

func (r *GoogleRepository) CountGoogleUserByEmail(ctx context.Context, email string) (total int, err error) {
	sql := r.Conn.WithContext(ctx).
		Table("users").
		Select("count(*) as total").
		Where("email = ?", email).
		Where("login_type = 'GOOGLE'").
		Scan(&total)
	if sql.Error != nil {
		return 0, sql.Error
	}
	return total, nil
}
