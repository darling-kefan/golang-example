package models

import (
	sqlEngine "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"log"
	"strconv"
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

func (p *publishsTable) PubsByIds(publishids interface{}) []map[string]string {
	if publishids == nil {
		return nil
	}
	pubids := make([]string, 0)
	switch v := publishids.(type) {
	case int:
		pubids = []string{strconv.Itoa(v)}
	case []int:
		for _, vv := range v {
			pubids = append(pubids, strconv.Itoa(vv))
		}
	case []string:
		pubids = v
	}
	
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
		rowCount++
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return res
}

// GetPubByPubid get publish by publish_id
func (p *publishsTable) JoinPubsByIds(publishids interface{}) []map[string]string {
	if publishids == nil {
		return nil
	}
	pubids := make([]string, 0)
	switch v := publishids.(type) {
	case int:
		pubids = []string{strconv.Itoa(v)}
	case []int:
		for _, vv := range v {
			pubids = append(pubids, strconv.Itoa(vv))
		}
	case []string:
		pubids = v
	}
	
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
		"cards.status as cardStatus",
		"cards.promote_type as promoteType",
		"cards.promote_route as promoteRoute",
		"advertisers.name as advertiser_name",
		"advertisers.tvmid as tvmid",
		"advertisers.shortcom as shortcom",
		"advertisers.logo as comlogo",
		"advertisers.industry",
		"advertisers.type as adverType",
		"advertisers.subtype as adverSubtype",
		"advertisers.istop as istop",
		"advertisers.version as versionNo",
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
func (p *publishsTable) LastestPubversion(pubid int, pubvers int) (ret map[string]interface{}) {
	var apiUrl string
	if pubvers == 2 {
		apiUrl = config.Get("ZTC_API_HOST_V2", "http://alpha.e.tvm.cn/restful/api/publish/pubversion") +
			"?pid=" + strconv.Itoa(pubid)
	} else {
		apiUrl = config.Get("ZTC_API_HOST", "http://alpha.e.tvm.cn/api/publish") +
			"?publishid=" + strconv.Itoa(pubid)
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
	if pubvers == 2 {
		for k, v := range m {
			switch k {
			case "version","target","time","deliverychannels":
				ret[k] = v
			case "keyword":
				kv := v.(map[string]interface{})
				ret["card"] = kv["card"]
			}
		}
	} else {
		if m["code"] != "0" {
			return
		}
		m = m["msg"].(map[string]interface{})
		for k, v := range m {
			switch k {
			case "version","target","time":
				ret[k] = v
			case "keyword":
				kv := v.(map[string]interface{})
				ret["card"] = kv["card"]
			}
		}
	}
	return
}
