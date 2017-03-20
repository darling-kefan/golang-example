package models

import (
	sqlEngine "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type cardsTable struct{}

func Cards() *cardsTable {
	return new(cardsTable)
}

func (c *cardsTable) CardsByIds(cardIds []string) []map[string]string {
	sql := "SELECT * FROM `cards` WHERE `deleted_at` IS NULL AND `id` IN (%s)"
	sql = fmt.Sprintf(sql, strings.Join(cardIds, ","))
	db, err := sqlEngine.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]sqlEngine.RawBytes, len(columns))
	scanVals := make([]interface{}, len(columns))
	for k := range values {
		scanVals[k] = &values[k]
	}

	rowCount := 0
	res := make([]map[string]string, len(cardIds))
	for rows.Next() {
		err = rows.Scan(scanVals...)
		if err != nil {
			log.Fatal(err)
		}

		var value string
		line := make(map[string]string, len(columns))
		for k, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			line[columns[k]] = value
		}

		res[rowCount] = line
		rowCount++
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}
