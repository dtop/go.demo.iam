package models

import (
	"database/sql"

	"github.com/dtop/go.demo.iam/iam/wrappers"
	"github.com/dtop/go.ginject"
	"gopkg.in/oauth2.v3"
)

type (

	// Client eventually represents one client
	Client struct {
		Ident   string
		Secret  string
		Scopes  []string
		Domains []string
	}

	// ClientStorage represents the actual storage handler
	ClientStorage struct {
		Db *wrappers.MySQL `inject:"db"`
	}
)

// #################### ClientStorage

// NewClientStorage creates and returns a new client storage including all dependencies
func NewClientStorage(dep ginject.Injector) (oauth2.ClientStore, error) {

	cs := &ClientStorage{}
	if err := dep.Apply(cs); err != nil {
		return nil, err
	}

	return cs, nil
}

// GetByID returns this object loaded with the given id (ClientStore Interface)
func (cs *ClientStorage) GetByID(id string) (oauth2.ClientInfo, error) {

	rows, err := cs.Db.Query(func(link *sql.DB) (*sql.Rows, error) {

		qry := `SELECT c.id, c.secret, cd.domain, cs.scope FROM clients c
				LEFT JOIN clients_domains cd ON c.id = cd.client_id
				LEFT JOIN clients_scopes cs ON c.id = cs.client_id
				WHERE c.id = ?`
		return link.Query(qry, id)
	})

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var (
		ident  string
		secret string
		domain string
		scope  string
	)

	cli := &Client{
		Scopes:  make([]string, 0),
		Domains: make([]string, 0),
	}

	for rows.Next() {

		if err := rows.Scan(&ident, &secret, &domain, &scope); err != nil {
			return nil, err
		}

		if domain != "" {
			cli.Domains = append(cli.Domains, domain)
		}

		if scope != "" {
			cli.Scopes = append(cli.Scopes, scope)
		}

	}

	cli.Ident = ident
	cli.Secret = secret

	return cli, nil
}

// ######################## Client

// GetID returns the client ID
func (cli *Client) GetID() string {
	return cli.Ident
}

// GetSecret returns the client secret
func (cli *Client) GetSecret() string {
	return cli.Secret
}

// GetDomain returns the first domain (mostly enaugh)
func (cli *Client) GetDomain() string {
	return cli.Domains[0]
}

// GetUserID returns the userid which we do not support
func (cli *Client) GetUserID() string {
	return ""
}
