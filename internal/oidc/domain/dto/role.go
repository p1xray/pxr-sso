package dto

// Role is a DTO with user role data.
type Role struct {
	id          int64
	code        string
	permissions []string
}

func NewRole(id int64, code string, permissions []string) Role {
	return Role{
		id:          id,
		code:        code,
		permissions: permissions,
	}
}

func (r *Role) ID() int64 {
	return r.id
}

func (r *Role) Code() string {
	return r.code
}

func (r *Role) Permissions() []string {
	return r.permissions
}
