package models

import (
	"fmt"
	"log"
	"time"
)

type advversions struct{}

func Advversions() *advversions {
	return new(advversions)
}

func (a *advversions) BasicAdverIds() (ret []int) {
	date := time.Now().Format("2006-01-02")
	sql  := "SELECT `advertisers`.`id` FROM `advertisers` " +
		"LEFT JOIN `advversions` ON `advertisers`.`tvmid` = `advversions`.`tvmid` " +
		"WHERE `advversions`.`status` = 1 " +
		"AND `advversions`.`stime` <= '%s' AND `advversions`.`etime` >= '%s'"
	sql  = fmt.Sprintf(sql, date, date)

	rows, err := Conn.Query(sql)
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
