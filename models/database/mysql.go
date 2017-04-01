/*package database

import (
	"fmt"
	"strconv"
	"database/sql"
)

type MySQLConfig struct {
	Username, Password string
	Host string
	Port int
}

func (cf *MySQLConfig) dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cf.Username, cf.Password, cf.Host, cf.Port, db)
}

type Conn *sql.DB

func NewConn(cf *MySQLConfig, db string) () {
	var err error
	Conn, err = sql.Open("mysql", cf.dsn(db))
	if err != nil {
		return err.Error())
	}
	if err := Conn.Ping(); err != nil {
		Conn.Close()
		panic(err.Error())
	}
}*/





