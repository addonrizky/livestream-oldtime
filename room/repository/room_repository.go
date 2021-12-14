package repository

import (
	"context"

	room "github.com/asumsi/livestream/room"
	"github.com/asumsi/livestream/room/models"
	"gorm.io/gorm"
)

type RoomRepository struct {
	Conn *gorm.DB
}

func NewRoomRepository(db *gorm.DB) room.RoomRepository {
	return &RoomRepository{
		Conn: db,
	}
}

func (r *RoomRepository) CountRooms(ctx context.Context) (total int, err error) {
	sql := r.Conn.WithContext(ctx).
		Table("rooms").
		Select("count(*) as total").
		Scan(&total)
	if sql.Error != nil {
		return 0, sql.Error
	}

	return total, nil
}

func (r *RoomRepository) GetRooms(ctx context.Context, limit, offset int) (res []models.Room, err error) {
	sql := r.Conn.WithContext(ctx).
		Table("rooms").
		Select("id, title, description, schedule_time, is_premium, created_by, created_date, updated_by, updated_date").
		Limit(limit).
		Offset(offset).
		Scan(&res)
	if sql.Error != nil {
		return res, sql.Error
	}

	return res, nil
}
