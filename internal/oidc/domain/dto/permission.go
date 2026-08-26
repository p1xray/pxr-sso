package dto

type Permission struct {
	id          int64
	code        string
	description string
}

func NewPermission(id int64, code, description string) Permission {
	return Permission{
		id:          id,
		code:        code,
		description: description,
	}
}

func (p *Permission) ID() int64 {
	return p.id
}

func (p *Permission) Code() string {
	return p.code
}

func (p *Permission) Description() string {
	return p.description
}
