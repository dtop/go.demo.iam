package wrappers

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	dsn  string
	link *sql.DB
}

// #################### MySQL

// NewMySQL creates a new mysql object
func NewMySQL(sock, host string, port int, user, pass, base string) *MySQL {

	dsnCreate := func(sock, host string, port int, user, pass, base string) string {

		var conn = ""

		if sock != "" {
			conn = fmt.Sprintf("unix(%v)", sock)
		}

		if host != "" {

			if port <= 0 {
				port = 3306
			}

			conn = fmt.Sprintf("tcp(%v:%v)", host, port)
		}

		return fmt.Sprintf("%v:%v@%v/%v?charset=utf8&loc=Local&parseTime=true", user, pass, conn, base)

	}

	return &MySQL{dsn: dsnCreate(sock, host, port, user, pass, base)}
}

// Open opens a new database connection
func (my *MySQL) Open() (link *sql.DB, err error) {

	link, err = sql.Open("mysql", my.dsn)
	if err == nil {
		my.link = link
	}
	return
}

// Close closes a database connection
func (my *MySQL) Close() {

	if my.link != nil {
		my.link.Close()
		my.link = nil
	}
}
