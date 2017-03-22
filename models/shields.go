package models

import (
	sqlEngine "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"

	"config"
)

type shieldsTable struct{}

func Shields() *shieldsTable {
	return new(shieldsTable)
}

// 查找屏蔽记录
// stype     屏蔽类型 1-分类，2-创意，3-广告主
// products  产品id，1直通车基础版 2-直通车专业版 3-KA
// markerIds 分类id/创意id/广告主id，与参数$type对应
// fields    字段
func (s *shieldsTable) Find(stype int, products []string, markerIds []string, fields []string) (ret []map[string]string) {
	sql := "SELECT " + strings.Join(fields, ",") + " FROM `shields` WHERE `status` = 1"
	if stype == 1 || stype == 2 || stype == 3 {
		sql += " AND `type` = %s"
	}
	if productIds != nil || len(productIds) != 0 {
		sql += " AND `product_id` in (%s)"
	}
	if markerIds != nil || len(markerIds) != 0 {
		sql += " AND `marker_id` in (%s)"
	}
	sql = fmt.Sprintf(sql, stype, strings.Join(products, ","), strings.Join(markerIds, ","))

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
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	rowCount := 0
	ret = make([]map[string]string)
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
				value = string(col)
			}
			line[columns[k]] = value
		}
		ret = append(ret, line)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
}
