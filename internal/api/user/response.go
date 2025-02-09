package user

type userSessionGetResponse struct {
	ID          string  `json:"id"`
	Nickname    *string `json:"nickname"`
	Permissions int32   `json:"permissions"`
	Cash        uint    `json:"cash"`
}

type userGetResponse struct {
	ID       uint64  `json:"id"`
	Nickname *string `json:"nickname"`
}
