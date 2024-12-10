package user

type UpdateUserBody struct {
	Nickname    *string `json:"nickname"`
	Permissions *int32  `json:"permissions"`
	Cash        *int32  `json:"cash"`
}
