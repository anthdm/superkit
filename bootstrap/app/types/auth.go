package types

// AuthUser represents an user that might be authenticated.
type AuthUser struct {
	ID       uint
	Email    string
	LoggedIn bool
}

// Check should return true if the user is authenticated.
// See handlers/auth.go.
func (user AuthUser) Check() bool {
	return user.ID > 0 && user.LoggedIn
}
