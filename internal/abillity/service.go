package abillity

import "github.com/dipress/crmifc/internal/user"

const (
	// ADMIN is a role name.
	ADMIN = "Admin"
)

// UserAbillity allows checking ability to view by user role.
type UserAbillity struct{}

// CanAdmin checks that user have a admin role.
func (a UserAbillity) CanAdmin(u *user.User) bool {
	if u.Role.Name == ADMIN {
		return true
	}
	return false
}
