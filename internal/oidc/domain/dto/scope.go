package dto

// Scope is a DTO with scope data.
type Scope struct {
	id          int64
	code        string
	name        string
	description string
}

func NewScope(id int64, code, name, description string) Scope {
	return Scope{
		id:          id,
		code:        code,
		name:        name,
		description: description,
	}
}

func (s *Scope) ID() int64 {
	return s.id
}

func (s *Scope) Code() string {
	return s.code
}

func (s *Scope) Name() string {
	return s.name
}

func (s *Scope) Description() string {
	return s.description
}
