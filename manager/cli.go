package manager

import (
	"github.com/ernestio/ernest-cli/model"

	eclient "github.com/ernestio/ernest-go-sdk/client"
	econfig "github.com/ernestio/ernest-go-sdk/config"
)

// Client : ...
type Client struct {
	cli          *eclient.Client
	cfg          *model.Config
	user         *User
	session      *Session
	notification *Notification
	role         *Role
}

// New : ...
func New(config *model.Config) *Client {
	client := eclient.New(
		econfig.New(config.URL).WithToken(config.Token),
	)
	return &Client{cli: client, cfg: config}
}

// User : User wrapper lazy load
func (c *Client) User() *User {
	if c.user == nil {
		c.user = &User{cli: c.cli}
	}
	return c.user
}

// Session : Session wrapper lazy load
func (c *Client) Session() *Session {
	if c.session == nil {
		c.session = &Session{cli: c.cli}
	}
	return c.session
}

// Notification : Notification wrapper lazy load
func (c *Client) Notification() *Notification {
	if c.notification == nil {
		c.notification = &Notification{cli: c.cli}
	}
	return c.notification
}

// Role : Role wrapper lazy load
func (c *Client) Role() *Role {
	if c.role == nil {
		c.role = &Role{cli: c.cli}
	}
	return c.role
}

// Cli : gets the internal eclient.Client
func (c *Client) Cli() *eclient.Client {
	return c.cli
}

// Config : ...
func (c *Client) Config() *model.Config {
	return c.cfg
}
