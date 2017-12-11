package manager

import (
	"github.com/ernestio/ernest-cli/model"

	eclient "github.com/ernestio/ernest-go-sdk/client"
	econfig "github.com/ernestio/ernest-go-sdk/config"
)

// Client : ...
type Client struct {
	cli     *eclient.Client
	cfg     *model.Config
	user    *user
	session *session
}

// New : ...
func New(config *model.Config) *Client {
	client := eclient.New(
		econfig.New(config.URL).WithToken(config.Token),
	)
	return &Client{cli: client, cfg: config}
}

// User ...
func (c *Client) User() *user {
	if c.user == nil {
		c.user = &user{cli: c.cli}
	}
	return c.user
}

// Session : ...
func (c *Client) Session() *session {
	if c.session == nil {
		c.session = &session{cli: c.cli}
	}
	return c.session
}

// Cli ...
func (c *Client) Cli() *eclient.Client {
	return c.cli
}

// Config : ...
func (c *Client) Config() *model.Config {
	return c.cfg
}
