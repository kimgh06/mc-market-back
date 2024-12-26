package permissions

type UserPermission = int

const (
	Administrator UserPermission = 1 << iota
	ListUsers
	ManageProducts
	ManageUsers
	ManageArticles
)
