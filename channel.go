package cicada

import (
	_ "fmt"
	_ "os"
	"log"
	"strings"
	_ "reflect"
	"database/sql"

	"github.com/urfave/cli"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
)

func ChannelAction(c *cli.Context) error {
	redisConn, err := redis.Dial("tcp", "10.10.51.14:6311")
	if err != nil {
		log.Fatalf("Redis connect failed: %v\n", err)
	}
	defer redisConn.Close()

	// 获取所有channel mids
	ret, err := redis.Values(redisConn.Do("lrange", "adview:channellist", 0, -1))
	chanMids := make([]string, len(ret))
	for k, v := range ret {
		chanMids[k] = string(v.([]byte))
	}

	// 获取mysql已存在的channel mids
	etvmConn, err := sql.Open("mysql", "root:tvmining@123@tcp(10.10.72.64:3306)/etvm")
	if err != nil {
		log.Fatalf("mysql: could not get a connection: %v\n", err)
	}
	defer etvmConn.Close()
	if err := etvmConn.Ping(); err != nil {
		etvmConn.Close()
		log.Fatalf("mysql: Could not establish a good connection: %v\n", err)
	}

	query := "SELECT `mid` FROM `channels` WHERE `id` >= 10000"
	stmt, err := etvmConn.Prepare(query)
	if err != nil {
		log.Fatalf("mysql: %v\n", err)
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Fatalf("mysql: %v\n", err)
	}
	defer rows.Close()

	var existMids []string
	for rows.Next() {
		var mid string
		if err := rows.Scan(&mid); err != nil {
			log.Fatalf("mysql: %v\n", err)
		}
		existMids = append(existMids, mid)
	}
	if err := rows.Err(); err != nil {
		log.Fatalf("mysql: %v\n", err)
	}

	var insertMids []string
	for _, v := range chanMids {
		isExist := false
		for _, vv := range existMids {
			if v == vv {
				isExist = true
			}
		}
		if !isExist {
			insertMids = append(insertMids, v)
		}
	}

	log.Println(chanMids, existMids, insertMids)
	
	if len(insertMids) == 0 {
		return nil
	}

	for _, v := range insertMids {
		key := "channel:mid:"+v
		res, err := redis.StringMap(redisConn.Do("HGETALL", key))
		if err != nil {
			log.Fatalf("redis: %v\n", err)
		}

		fields := []string{
			"`id`",
			"`title`",
			"`mid`",
			"`sn`",
			"`subtitle`",
			"`wx_token`",
			"`yyyappid`",
			"`type`",
			"`parent_id`",
			"`desc`",
			"`logo`",
			"`heat`",
			"`online`",
			"`is_owner`",
			"`created_at`",
			"`updated_at`",
		}
		query := "INSERT INTO `channels`("+strings.Join(fields, ",")+") VALUES(?"+strings.Repeat(",?",len(fields)-1)+")"
		stmt, err := etvmConn.Prepare(query)
		if err != nil {
			log.Fatalf("mysql: %v\n", err)
		}
		defer stmt.Close()
		args := []interface{} {
			res["id"],
			res["title"],
			res["mid"],
			res["sn"],
			res["subtitle"],
			res["wx_token"],
			res["yyyappid"],
			res["type"],
			res["parent_id"],
			"adview",
			res["logo"],
			res["heat"],
			res["online"],
			res["is_owner"],
			res["created_at"],
			res["updated_at"],
		}
		r, err := stmt.Exec(args...)
		if err != nil {
			log.Fatalf("mysql: could not execute statement: %v", err)
		}
		if _, err := r.RowsAffected(); err != nil {
			log.Fatalf("mysql: could not get rows affected: %v", err)
		}
		lastInsertId, _ := r.LastInsertId()
		log.Printf("%v\n", lastInsertId)
	}
	
	return nil
}
