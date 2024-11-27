package user

type userSessionGetResponse struct {
	ID          uint64  `json:"id"`
	Nickname    *string `json:"nickname"`
	Permissions int32   `json:"permissions"`
}

type userGetResponse struct {
	ID       uint64  `json:"id"`
	Nickname *string `json:"nickname"`
}
