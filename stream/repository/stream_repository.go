package repository

import (
	"context"

	"github.com/asumsi/livestream/stream/models"
	"gorm.io/gorm"
)

type StreamRepository struct {
	Conn *gorm.DB
}

func NewStreamRepository(db *gorm.DB) *StreamRepository {
	return &StreamRepository{
		Conn: db,
	}
}

func (r *StreamRepository) GetStreams(ctx context.Context, id int64) (res []*models.Stream, err error) {
	return nil, nil
}
