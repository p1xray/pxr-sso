package dto

// Client is a DTO with client data.
type Client struct {
	id           int64
	code         string
	secretKey    string
	audiences    []string
	redirectURI  []string
	scope        []string
	defaultRoles []Role
}

func NewClient(
	id int64,
	code,
	secretKey string,
	audiences,
	redirectURI,
	scope []string,
	defaultRoles []Role,
) Client {
	return Client{
		id:           id,
		code:         code,
		secretKey:    secretKey,
		audiences:    audiences,
		redirectURI:  redirectURI,
		scope:        scope,
		defaultRoles: defaultRoles,
	}
}

func (c *Client) ID() int64 {
	return c.id
}

func (c *Client) Code() string {
	return c.code
}

func (c *Client) SecretKey() string {
	return c.secretKey
}

func (c *Client) Audiences() []string {
	return c.audiences
}

func (c *Client) RedirectURI() []string {
	return c.redirectURI
}

func (c *Client) Scope() []string {
	return c.scope
}

func (c *Client) DefaultRoles() []Role {
	return c.defaultRoles
}
