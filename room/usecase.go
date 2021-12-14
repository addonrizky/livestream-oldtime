package room

import (
	"context"

	"github.com/asumsi/livestream/room/models"
)

type RoomUsecase interface {
	GetRooms(ctx context.Context, limit, page int) (res models.RoomAttribute, err error)
}
