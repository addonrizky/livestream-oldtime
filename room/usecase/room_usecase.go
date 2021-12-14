package usecase

import (
	"context"
	"time"

	room "github.com/asumsi/livestream/room"
	"github.com/asumsi/livestream/room/models"
	"github.com/sirupsen/logrus"
)

type roomUsecase struct {
	roomRepo       room.RoomRepository
	contextTimeout time.Duration
}

func NewRoomUsecase(a room.RoomRepository, timeout time.Duration) room.RoomUsecase {
	return &roomUsecase{
		roomRepo:       a,
		contextTimeout: timeout,
	}
}

func (a *roomUsecase) GetRooms(ctx context.Context, limit, page int) (res models.RoomAttribute, err error) {
	ctx, cancel := context.WithTimeout(ctx, a.contextTimeout)
	defer cancel()

	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	countRoom, err := a.roomRepo.CountRooms(ctx)
	if err != nil {
		logrus.Infof("CountRooms %v", err)
		return res, err
	}

	roomList, err := a.roomRepo.GetRooms(ctx, limit, offset)
	if err != nil {
		logrus.Infof("GetRooms %v", err)
		return res, err
	}

	res.Page = page
	res.DataPerPage = limit
	res.TotalData = countRoom
	res.Data = roomList

	return res, nil
}
