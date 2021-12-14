package models

type InitSessionRequest struct {
	RoomId      string `json:"room_id" validate:"required" example:"9918"`
	UserId 		string `json:"user_id" validate:"required" example:"adhon.rizky@gmail.com"`
	OfferSession string `json:"offer_session" validate:"required" example:"nmhbf93y21592r3fhjgoewhgveb28356932kffkl"`
}

type StartSessionRequest struct {
	RoomId      string `json:"room_id" validate:"required"`
}

type EndSessionRequest struct {
	PID      string `json:"pid" validate:"required"`
}