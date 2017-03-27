package districts

import (
	"fmt"
	"errors"
	"strconv"
	"database/sql"
	"github.com/garyburd/redigo/redis"
)

type District struct {
	Id      string
	Level   int
	Members []string
}

func (ds *District) ToRedis() error {
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		return err
	}
	defer c.Close()

	var rkey string
	switch ds.Level {
	case 1:
		rkey = "districts:country:1"
	case 2:
		rkey = "districts:province:" + ds.Id
	case 3:
		rkey = "districts:city:" + ds.Id
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
	
	return nil
}

func Source() ([]*District, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", "root", "123456", "localhost", 3306, "etvm")
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
	for rows.Next() {
		err := rows.Scan(&id, &did, &pid, &level)
		if err != nil {
			return nil, err
		}

		ds := new(District)
		if level < 3 {
			ds.Id = id
			ds.Level = level
		}
		switch level {
		case 0:
			ds.Id = id
			ds.Level = level
		case 1:
			ds.Id = id
			ds.Level = level
		case 2:
			ds.Id = id
			ds.Level = level
		}
	}
	
	return districts, nil
}

