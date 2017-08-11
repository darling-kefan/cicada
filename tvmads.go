package cicada

import (
	"fmt"
	"os"
	"log"
	"time"
	"math"
	"bytes"
	"runtime"
	"strings"
	"strconv"
	_ "reflect"
	"database/sql"

	"github.com/urfave/cli"
	_ "github.com/go-sql-driver/mysql"
	"github.com/garyburd/redigo/redis"
)

func UniformSpeedAction(ctx *cli.Context) error {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)

	env, segment, day := ctx.GlobalString("env"), ctx.Int("segment"), ctx.Int("day")

	var (
		redisHost string
		redisPort string
		ejobHost  string
		ejobPort  string
		ejobUser  string
		ejobPass  string
	)
	switch env {
	case "local":
		redisHost = "127.0.0.1"
		redisPort = "6379"
		ejobHost  = "127.0.0.1"
		ejobPort  = "3306"
		ejobUser  = "root"
		ejobPass  = "123456"
	case "dev","develop":
		redisHost = "10.10.51.14"
		redisPort = "6311"
		ejobHost  = "10.10.72.64"
		ejobPort  = "3306"
		ejobUser  = "root"
		ejobPass  = "tvmining@123"
	case "prod","production":
		redisHost = "127.0.0.1"
		redisPort = "6311"
		ejobHost  = "10.66.189.75"
		ejobPort  = "3306"
		ejobUser  = "ejob"
		ejobPass  = "CcN28T-V4#T!"
	default:
		log.Fatalf("[error] parameter env is not valid, %v", env)
		return nil
	}

	// Get publish IDs
	pubids := make([]int64, ctx.NArg())
	for i:=0; i<ctx.NArg(); i++ {
		pubid, err := strconv.ParseInt(ctx.Args().Get(i), 10, 64)
		if err != nil {
			_, fn, line, _ := runtime.Caller(1)
			log.Fatalf("[error] %s:%d %v", fn, line, err)
		}
		pubids[i] = pubid
	}
	// Get ads
	dsn := fmt.Sprintf("%s:%s@tcp([%s]:%s)/ejob", ejobUser, ejobPass, ejobHost, ejobPort)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", fn, line, err)
	}
	if err := conn.Ping(); err != nil {
		conn.Close()
		_, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", fn, line, err)
	}

	uniformAds := uniformSpeedAds(pubids, day, conn)
	if len(uniformAds) == 0 {
		return nil
	}

	rc, err := redis.Dial("tcp", redisHost+":"+redisPort)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, ad := range uniformAds {
		if err := toredis10(ad, segment, day, rc); err != nil {
			_, fn, line, _ := runtime.Caller(1)
			log.Fatalf("[error] %s:%d %v", fn, line, err)
		}
	}

	return nil
}

func uniformSpeedAds(pubids []int64, day int, conn *sql.DB) []*UniformAd {
	tdate := time.Now().AddDate(0, 0, day)
	sql := "SELECT `publish_id`,`t_hours`,`s_bid`,`s_ceiling` FROM `pubcaches`"
	sqlArgs := make([]interface{}, len(pubids)+1)
	sqlArgs[0] = tdate.Format("2006-01-02")
	if len(pubids) > 0 {
		sql = sql+" WHERE `speed_type` = 1 AND `t_date` = ? AND `publish_id` IN (?"+strings.Repeat(",?",len(pubids)-1)+")"
		for k, v := range pubids {
			sqlArgs[k+1] = v
		}
	} else {
		sql = sql+" WHERE `speed_type` = 1 AND `t_date` = ?"
	}
	stmt, err := conn.Prepare(sql)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", fn, line, err)
	}
	defer stmt.Close()
	rows, err := stmt.Query(sqlArgs...)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Fatalf("[error] %s:%d %v", fn, line, err)
	}
	defer rows.Close()

	var uniformAds []*UniformAd
	for rows.Next() {
		var (
			id      int64
			bid     int64
			ceiling int64
			hours   string
		)
		if err := rows.Scan(&id, &hours, &bid, &ceiling); err != nil {
			_, fn, line, _ := runtime.Caller(1)
			log.Fatalf("[error] %s:%d %v", fn, line, err)
		}
		uniformAd := &UniformAd{
			Id:      id,
			Bid:     bid,
			Ceiling: ceiling,
			Hours:   strings.Split(hours, ","),
		}
		uniformAds = append(uniformAds, uniformAd)
	}

	return uniformAds
}

type UniformAd struct {
	Id      int64
	Bid     int64
	Ceiling int64
	Hours   []string
}

func toredis10(ad *UniformAd, segment int, day int, rc redis.Conn) error {
	now := time.Now()

	//标记当前同步时间
	//markTime := now.Add(10 * time.Minute)
	markTime := now
	if day != 0 {
		year, month, day := now.AddDate(0, 0, day).Date()
		loc, _ := time.LoadLocation("Asia/Shanghai")
		markTime = time.Date(year, month, day, 0, 0, 0, 0, loc)
	}

	// 已消耗金额
	budgetKey := "budget:pub:"+strconv.FormatInt(ad.Id, 10)+":"+markTime.Format("20060102")
	consume, _ := redis.Int64(rc.Do("HGET", budgetKey, "consume"))

	// 计算总余额
	balance := ad.Ceiling - consume
	// 计算总曝光数
	shows := int64(math.Ceil(float64(balance)/float64(ad.Bid)))
	
	// 获取每小时曝光基数
	base := exposureBase(shows, ad.Hours, day)

	log.Printf("balance: %d, bid: %d, shows: %d, base: %d\n", balance, ad.Bid, shows, base)

	var buf bytes.Buffer
	var sendCount int
	for i:=0; i<=segment; i++ {
		tt := markTime.Add(time.Duration(i*10)*time.Minute)
		ts := tt.Format("1504")

		buf.WriteString("ys:")
		buf.WriteString(ts[0:3])
		buf.WriteString("0:")
		buf.WriteString(strconv.FormatInt(ad.Id, 10))
		rkey := buf.String()
		buf.Reset()

		// 判断是否更新当前限速时段
		if i == 0 {
			ok, err := redis.Bool(rc.Do("EXISTS", rkey))
			if err != nil {
				return err
			}
			if ok {
				continue
			}
		}
		
		// tt小时总曝光量
		theShows := base
		if isHotHour(tt.Hour()) {
			theShows = int64(math.Ceil(1.2 * float64(base)))
		}

		// 每10分钟平均曝光量
		theShowPart := int64(math.Ceil(float64(theShows) / 6))

		//Send writes the command to the connection's output buffer.
		rc.Send("SET", rkey, theShowPart)
		rc.Send("EXPIRE", rkey, 2*3600)
		log.Printf("redis: Send SET %s %v\n", rkey, theShowPart)
		log.Printf("redis: Send EXPIRE %s %v\n", rkey, 2*3600)

		sendCount = sendCount + 2
	}

	//Flush flushes the connection's output buffer to the server
	rc.Flush()
	for i:=0; i<sendCount; i++ {
		//Receive reads a single reply from the server.
		if r, err := rc.Receive(); err == nil {
			log.Printf("redis: Receive: %v", r)
		} else {
			_, fn, line, _ := runtime.Caller(1)
			log.Fatalf("[error] %s:%d %v", fn, line, err)
		}
	}

	//Check publishs:id:<id> -> uniform_speed existed
	pubKey := fmt.Sprintf("publishs:id:%d", ad.Id)
	if ok, err := redis.Bool(rc.Do("HEXISTS", pubKey, "uniform_speed")); err == nil && !ok {
		if ret, err := rc.Do("HSET", pubKey, "uniform_speed", 1); err == nil {
			log.Printf("redis: HSET %s uniform_speed 1 -> %v", pubKey, ret)
		}
	}
	
	return nil
}

// 计算小时曝光基数
// 
// 计算公式：
// 热门曝光时段个数:m, 非热门曝光时段个数:n, 曝光基数:x
// m*1.2*x + n*x = balance
func exposureBase(exposures int64, hours []string, day int) int64 {
	m, n, now := 0, 0, time.Now()

	//标记小时数，用于过滤逝去小时。默认0, 即非今天标记小时数取0。
	markHour := 0
	if day == 0 {
		markHour = now.Hour()
		// 如果当前时间的分钟数大于等于50,则同步下一小时数据;否则取当前小时数。同时必须保证下一个小时必须在今天。
		if now.Minute() >= 50 && now.Hour() < 23 {
			markHour = now.Add(1*time.Hour).Hour() + 1
		}
	}

	for _, v := range hours {
		var hour int
		if v != "00" {
			var err error
			if hour, err = strconv.Atoi(strings.TrimLeft(v, "0")); err != nil {
				_, fn, line, _ := runtime.Caller(1)
				log.Fatalf("[error] %s:%d %v", fn, line, err)
			}
		}
		if hour < markHour {
			continue
		}
		for _, hotHour := range HotHours {
			if hotHour == hour {
				m++
			} else {
				n++
			}
		}
	}

	base := math.Ceil(float64(exposures)/(float64(m)*1.2 + float64(n)))
	return int64(base)
}

// 11,12,18,19,20,21,22称为热门曝光时段，热门曝光时段曝光比普通曝光时段多20%
var HotHours = []int{11, 12, 18, 19, 20, 21, 22}

// 判断是否为热门小时
func isHotHour(hour int) bool {
	for _, v := range HotHours {
		if v == hour {
			return true
		}
	}
	return false
}
