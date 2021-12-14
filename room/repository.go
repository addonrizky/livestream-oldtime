package room

import (
	"context"

	"github.com/asumsi/livestream/room/models"
)

type RoomRepository interface {
	CountRooms(ctx context.Context) (total int, err error)
	GetRooms(ctx context.Context, limit, offset int) (res []models.Room, err error)
}
