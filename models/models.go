package models

import (
	"fmt"

	"config"
)

var dsn, driver string

func init() {
	host := config.Get("DB_HOST", "127.0.0.1")
	port := config.Get("DB_PORT", "3306")
	database := config.Get("DB_DATABASE", "etvm")
	username := config.Get("DB_USERNAME", "root")
	password := config.Get("DB_PASSWORD", "123456")
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	driver = config.Get("DB_CONNECTION", "mysql")
}
