package models

import (
	sqlEngine "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"time"
)

type advversionsTable struct{}

func Advversions() *advversionsTable {
	return new(advversionsTable)
}

func (a *advversionsTable) BasicAdverIds() (ret []int) {
	date := time.Now().Format("2006-01-02")
	sql  := "SELECT `advertisers`.`id` FROM `advertisers` " +
		"LEFT JOIN `advversions` ON `advertisers`.`tvmid` = `advversions`.`tvmid` " +
		"WHERE `advversions`.`status` = 1 " +
		"AND `advversions`.`stime` <= '%s' AND `advversions`.`etime` >= '%s'"
	sql  = fmt.Sprintf(sql, date, date)
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
	var advid int
	for rows.Next() {
		err = rows.Scan(&advid)
		if err != nil {
			log.Fatal(err)
		}
		ret = append(ret, advid)
	}
	return
}
