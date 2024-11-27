package responses

import (
	"maple/internal/schema"
	"time"
)

type FullUser struct {
	ID          uint64    `json:"id"`
	Nickname    *string   `json:"nickname"`
	Permissions int32     `json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ShortUser struct {
	ID          uint64  `json:"id"`
	Nickname    *string `json:"username"`
	Permissions int32   `json:"permissions"`
}

func (u *FullUser) ToShortUser() ShortUser {
	return ShortUser{
		ID:          u.ID,
		Nickname:    u.Nickname,
		Permissions: u.Permissions,
	}
}

func FullUserFromSchema(user *schema.User) FullUser {
	var converted FullUser
	converted.ID = uint64(user.ID)
	converted.CreatedAt = user.CreatedAt
	converted.UpdatedAt = user.UpdatedAt
	converted.Permissions = user.Permissions
	if user.Nickname.Valid {
		converted.Nickname = &user.Nickname.String
	}

	return converted
}
