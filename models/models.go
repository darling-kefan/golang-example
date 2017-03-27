package models

import (
	"fmt"
	"strconv"
	"database/sql"
	"config"
)

var (
	Conn *sql.DB
)

func init() {
	port, _ := strconv.Atoi(config.Get("DB_PORT", "3306"))
	newConn(&MySQLConfig{
		Username: config.Get("DB_USERNAME", "root"),
		Password: config.Get("DB_PASSWORD", "123456"),
		Host: config.Get("DB_HOST", "127.0.0.1"),
		Port: port,
	}, "etvm")
}

type MySQLConfig struct {
	Username, Password string
	Host string
	Port int
}

func (cf *MySQLConfig) dsn(db string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cf.Username, cf.Password, cf.Host, cf.Port, db)
}

func newConn(cf *MySQLConfig, db string) {
	var err error
	// singleton
	if Conn == nil {
		Conn, err = sql.Open("mysql", cf.dsn(db))
		if err != nil {
			panic(err.Error())
		}
	}

	if err := Conn.Ping(); err != nil {
		Conn.Close()
		panic(err.Error())
	}
}

func Close() {
	Conn.Close()
}
