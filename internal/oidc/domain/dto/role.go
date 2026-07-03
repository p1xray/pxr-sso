package dto

// Role is a DTO with user role data.
type Role struct {
	id          int64
	code        string
	name        string
	description string
	permissions []Permission
}

func NewRole(id int64, code, name, description string, permissions []Permission) Role {
	return Role{
		id:          id,
		code:        code,
		name:        name,
		description: description,
		permissions: permissions,
	}
}

func (r *Role) ID() int64 {
	return r.id
}

func (r *Role) Code() string {
	return r.code
}

func (r *Role) Name() string {
	return r.name
}

func (r *Role) Description() string {
	return r.description
}

func (r *Role) Permissions() []Permission {
	return r.permissions
}
