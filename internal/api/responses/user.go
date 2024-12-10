package responses

import (
	"maple/internal/schema"
	"strconv"
	"time"
)

type FullUser struct {
	ID          string    `json:"id"`
	Nickname    *string   `json:"nickname"`
	Permissions int32     `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Cash        int       `json:"cash"`
}

type ShortUser struct {
	ID          string  `json:"id"`
	Username    *string `json:"username"`
	Nickname    *string `json:"nickname"`
	Permissions int32   `json:"permissions"`
	Cash        int     `json:"cash"`
}

func (u *FullUser) ToShortUser() ShortUser {
	return ShortUser{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Permissions: u.Permissions,
		Cash:        u.Cash,
	}
}

func FullUserFromSchema(user *schema.User) FullUser {
	var converted FullUser
	converted.ID = strconv.FormatUint(uint64(user.ID), 10)
	converted.CreatedAt = user.CreatedAt
	converted.UpdatedAt = user.UpdatedAt
	converted.Permissions = user.Permissions
	converted.Cash = int(user.Cash)
	if user.Nickname.Valid {
		converted.Nickname = &user.Nickname.String
	}

	return converted
}
