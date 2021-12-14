package modelgeneral

type Response struct {
	Code string      `json:"code"`
	ID   string      `json:"response_id"`
	Desc string      `json:"message"`
	Data interface{} `json:"data"`
}

func NewResponse(id string) *Response {
	return &Response{
		ID:   id,
		Code: "XX",
		Desc: "General Error",
		Data: new(struct{}),
	}
}

type LoginAttribute struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginResponse struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	Result  LoginAttribute `json:"result"`
}
