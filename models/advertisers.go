package models

import (
	sqlEngine "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

type advertisersTable struct{}

func Advertisers() *advertisersTable {
	return new(advertisersTable)
}

func (a *advertisersTable) AdversByIds(advertiserIds []string, fields interface{}) []map[string]string {
	if advertiserIds == nil || len(advertiserIds) == 0 {
		return nil
	}

	var fieldsStr string
	switch v := fields.(type) {
	case string:
		fieldsStr = v
	case []string:
		fieldsStr = strings.Join(v, ",")
	default:
		fieldsStr = "*"
	}
	sql := "SELECT %s FROM `advertisers` WHERE `id` IN (%s)"
	sql = fmt.Sprintf(sql, fieldsStr, strings.Join(advertiserIds, ","))

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

	var res []map[string]string
	for rows.Next() {
		err := rows.Scan(scanVals...)
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
		res = append(res, line)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}

func (a *advertisersTable) Qualifis(aids []string) []map[string]string {
	if aids == nil || len(aids) == 0 {
		return nil
	}

	fields := []string{
		"advertisers.*",
		"qualifis.id as qualifis_id",
		"qualifis.type as qualifis_type",
		"qualifis.on_type as qualifis_on_type",
		"qualifis.code as qualifis_code",
		"qualifis.codeimg as qualifis_codeimg",
		"qualifis.status as qualifis_status",
		"qualifis.msg as qualifis_msg",
		"qualifis.verifytime as qualifis_verifytime",
		"qualifis.passtime as qualifis_passtime",
		"qualifis.expiry as qualifis_expiry",
		"qualifis.created_at as qualifis_created_at",
		"qualifis.updated_at as qualifis_updated_at",
	}
	sql := "SELECT %s FROM `advertisers` LEFT JOIN `qualifis` ON `advertisers`.`id` = `qualifis`.`advertiser_id` " +
		"WHERE `advertisers`.`id` IN (%s)"
	sql = fmt.Sprintf(sql, strings.Join(fields, ","), strings.Join(aids, ","))
	fmt.Println(sql)

	db, err := sqlEngine.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(sql)
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	values := make([]sqlEngine.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for k := range values {
		scanArgs[k] = &values[k]
	}

	var res []map[string]string
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Fatal(err)
		}

		var value string
		line := make(map[string]string, len(columns))
		for k, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(value)
			}
			line[columns[k]] = value
		}
		res = append(res, line)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}
