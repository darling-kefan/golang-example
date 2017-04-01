package districts

import (
	"fmt"
	"errors"
	"strconv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
)

type District struct {
	Id      int
	Did     int
	Level   int
	Members []int
}

func ToRedis(districts []*District) error {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	defer c.Close()

	for _, ds := range districts {
		if len(ds.Members) == 0 {
			continue
		}
		var rkey string
		switch ds.Level {
		case 0:
			rkey = "districts:country:1"
		case 1:
			rkey = "districts:province:" + strconv.Itoa(ds.Id)
		case 2:
			rkey = "districts:city:" + strconv.Itoa(ds.Id)
		default:
			return errors.New("level is valid, level is " + strconv.Itoa(ds.Level))
		}

		// If you have slice of members, then copy key and members to []interface{}
		/*args := []interface{}{rkey}
	    for _, v := range members {
		    args = append(args, v)
	    }
	    if _, err := c.Do("SADD", args...); err != nil {
		    return err
	    }*/

		// redis.Args helper does same as previous
		if _, err := c.Do("SADD", redis.Args{}.Add(rkey).AddFlat(ds.Members)...); err != nil {
			return err
		}
	}
	
	return nil
}

func Source() ([]*District, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "tvmining@123", "10.10.72.64", 3306, "etvm")
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}
	listPre, err := conn.Prepare("SELECT `id`, `did`, `pid`, `leveltype` AS `level` FROM `districts` ORDER BY `id` ASC")
	if err != nil {
		return nil, err
	}
	rows, err := listPre.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var districts []*District
	var id, did, pid, level int
	dismap := make(map[int][]int, 0)
	for rows.Next() {
		err := rows.Scan(&id, &did, &pid, &level)
		if err != nil {
			return nil, err
		}
		switch level {
		case 0:
			dismap[did] = []int{}
			districts = append(districts, &District{Id: id, Did: did, Level: level})
		case 1,2:
			dismap[pid] = append(dismap[pid], id)
			districts = append(districts, &District{Id: id, Did: did, Level: level})
		case 3:
			dismap[pid] = append(dismap[pid], id)
		}
	}
	for _, v := range districts {
		v.Members = dismap[v.Did]
	}
	return districts, nil
}

