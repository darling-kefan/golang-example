package models

import (
	sqlEngine "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"config"
)

type publishsTable struct{}

func Publishs() *publishsTable {
	return new(publishsTable)
}

func (p *publishsTable) PubsByIds(pubids []string) []map[string]string {
	sql := "SELECT * FROM `publishs` WHERE `id` IN (%s)"
	db, err := sqlEngine.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query(fmt.Sprintf(sql, strings.Join(pubids, ",")))
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
	res := make([]map[string]string, len(pubids))
	for rows.Next() {
		fmt.Println(rowCount)
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error())
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
		fmt.Println(line)
		rowCount++
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}

// GetPubByPubid get publish by publish_id
func (p *publishsTable) JoinPubsByIds(pubids []string) []map[string]string {
	fields := []string{
		"publishs.id",
		"publishs.advertiser_id",
		"publishs.adposition_id",
		"publishs.title",
		"publishs.adplan_id",
		"publishs.stime",
		"publishs.etime",
		"publishs.bid",
		"publishs.ceiling",
		"publishs.allceiling",
		"publishs.status",
		"publishs.check",
		"publishs.billing",
		"publishs.mode",
		"publishs.type",
		"publishs.card_id",
		"publishs.group_id",
		"publishs.deleted_at",
		"publishs.updated_at",
		"publishs.created_at",
		"publishs.again_num",
		"cards.status as cards_status",
		"cards.promote_type as promote_type",
		"cards.promote_route as promote_route",
		"advertisers.shortcom as shortcom",
		"advertisers.logo as comlogo",
		"advertisers.industry",
		"advertisers.type as advertiser_type",
		"advertisers.subtype as advertiser_subtype",
		"advertisers.istop as istop",
		"advertisers.version as publish_version",
	}
	fieldsStr := strings.Join(fields, ",")

	sql := "SELECT " + fieldsStr + " FROM `publishs` " +
		"LEFT JOIN `cards` ON `publishs`.`card_id` = `cards`.`id` " +
		"LEFT JOIN `advertisers` ON `publishs`.`advertiser_id` = `advertisers`.`id` " +
		"WHERE `publishs`.`id` IN (%s)"

	db, err := sqlEngine.Open(driver, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Prepare action can not support returning multi rows ????
	/*stmtout, err := db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmtout.Close()
	rows, err := stmtout.Query(strings.Join(pubids, ","))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()*/

	rows, err := db.Query(fmt.Sprintf(sql, strings.Join(pubids, ",")))
	if err != nil {
		log.Fatal(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	// slice has sort
	values := make([]sqlEngine.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
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

// Get lastest pubversion
func (p *publishsTable) LastestPubversion(pubid string, pubvers int) (ret map[string]interface{}) {
	var apiUrl string
	if pubvers == 2 {
		apiUrl = config.Get("ZTC_API_HOST_V2", "http://alpha.e.tvm.cn/restful/api/publish/pubversion") +
			"?pid=" + pubid
	} else {
		apiUrl = config.Get("ZTC_API_HOST", "http://alpha.e.tvm.cn/api/publish") +
			"?publishid=" + pubid
	}
	res, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var f interface{}
	err = json.Unmarshal(data, &f)
	m := f.(map[string]interface{})
	ret = make(map[string]interface{})
	for k, v := range m {
		switch k {
		case "target","time","deliverychannels","version":
			ret[k] = v
		case "keyword":
			kv := v.(map[string]interface{})
			ret["card"] = kv["card"]
		}
	}
	return ret
}
