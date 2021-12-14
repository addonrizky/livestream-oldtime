package models

import "time"

type Room struct {
	ID           int64      `json:"id" gorm:"id"`
	Title        string     `json:"title" gorm:"title"`
	Description  string     `json:"description" gorm:"description"`
	ScheduleTime time.Time  `json:"schedule_time" gorm:"schedule_time"`
	IsPremium    bool       `json:"is_premium" gorm:"is_premium"`
	CreatedBy    int64      `json:"created_by" gorm:"created_by"`
	CreatedDate  time.Time  `json:"created_date" gorm:"created_date"`
	UpdatedBy    *int64     `json:"updated_by" gorm:"updated_by"`
	UpdatedDate  *time.Time `json:"updated_date" gorm:"updated_date"`
}

type RoomAttribute struct {
	Page        int    `json:"page"`
	DataPerPage int    `json:"dataPerPage"`
	TotalData   int    `json:"totalData"`
	Data        []Room `json:"data"`
}

type RoomResponse struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Result  RoomAttribute `json:"result"`
}
