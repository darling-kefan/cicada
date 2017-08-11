package cicada

import (
	"fmt"
	"os"
	"io"
	"time"
	"bufio"
	"strings"
	"encoding/json"

	"github.com/urfave/cli"
	"github.com/influxdata/influxdb/client/v2"
)

// batch size
const batchSize = 1000

func DetectAction(ctx *cli.Context) error {
	// detect type
	dt, source := ctx.String("type"), ctx.String("source")
	if source == "" {
		return cli.NewExitError("Error: The flag --source is required.", 1)
	}

	pids := []string{
		"c10de0d4ddbe24414cf0d0e3ab4ab6c95f89c74a-3717507",
		"ddff0693dc97c98ccd54b0d53e2fcfe3fc410564-3720595",
	}
	
	// Step 1: gather Gather data from source file and insert into influxdb
	fr, err := os.Open(source)
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: %s is not existed.", source), 2)
	}
	defer fr.Close()

	// Create a new InfluxDB HTTPClient
	ic, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: fmt.Sprintf("http://%s:8086", conf.InfluxDB.Server),
		Username: conf.InfluxDB.Username,
		Password: conf.InfluxDB.Password,
	})
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: %s", err.Error()), 3)
	}
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "detect_sync",
		Precision: "s",
	})
	if err != nil {
		return err
	}

	var parser Parser
	switch dt {
	case "ybt":
		parser = &Ztc{fr, ic, bp, "ybt"}
	case "tz":
		parser = &Ztc{fr, ic, bp, "tz"}
	case "cp":
		parser = &Ztc{fr, ic, bp, "cp"}
	case "ka":
		parser = &Ka{fr, ic, bp}
	default:
		return cli.NewExitError(fmt.Sprintf("Error: %s is invalid type.", dt), 4)
	}

	//if err := parser.Parse(); err != nil {
	//	return cli.NewExitError(fmt.Sprintf("Error: %s", err.Error()), 5)
	//}
	fmt.Println(parser)

	// Step 2: gather Gather data from source file and insert into influxdb
	if err := DisplayZtcData(pids, dt, ic); err != nil {
		return cli.NewExitError(fmt.Sprintf("Error: %s", err.Error()), 8)
	}
	
	return nil
}

type ZtcElement struct{
	Key     string
	OldPid  string
	OldVal  []string
	NewPid  string
	NewVal  []string
}

type ZtcElementList []ZtcElement

func (zel ZtcElementList) TableShow() error {
	return nil
}

func DisplayZtcData(pids []string, dt string, ic client.Client) error {
	query := fmt.Sprintf("SELECT * FROM ztc WHERE type = '%s' and %s ORDER BY time ASC", dt,
		fmt.Sprintf("pid_flag = '%s' or pid_flag = '%s'", pids[0], pids[1]))
	res, err := queryDB(ic, query)
	if err != nil {
		return err
	}
	//判断两个pid的先后顺序
	fmt.Println(res)
	
	return nil
}

func queryDB(cInt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: "detect_sync",
	}
	if response, err := cInt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

type Parser interface {
	Parse() error
}

type Ztc struct{
	Source io.Reader
	Ic     client.Client
	Bp     client.BatchPoints
	Type   string
}

func (ztc *Ztc) Parse() error {
	scanner := bufio.NewScanner(ztc.Source)
	var i int
	for scanner.Scan() {
		rowFields := strings.Fields(scanner.Text())
		if rowFields[3] == "benchmark:" || rowFields[3] == "The" ||
			rowFields[4] == "zrem" {
			continue
		}

		// Create a point and add to batch
		tags := map[string]string{
			"pid_flag": rowFields[2][1:len(rowFields[2])-1],
			"key": rowFields[3][1:len(rowFields[3])-1],
			"type": ztc.Type,
		}

		var value map[string]int64
		if err := json.Unmarshal([]byte(rowFields[5]), &value); err != nil {
			return err
		}
		pubIds := make([]string, len(value))
		for k, v := range value {
			kk := strings.Split(k, "-")
			pubIds[v-1] = kk[len(kk)-1]
		}
		fields := map[string]interface{}{
			//"key": rowFields[3][1:len(rowFields[3])-1],
			"value": strings.Join(pubIds, ","),
		}
		ts := rowFields[0][1:] + " " + rowFields[1][0:len(rowFields[1])-1]
		ti, err := time.Parse("2006-01-02 15:04:05", ts)
		if err != nil {
			return err
		}
		pt, err := client.NewPoint(
			"ztc",
			tags,
			fields,
			ti,
		)
		if err != nil {
			return err
		}
		ztc.Bp.AddPoint(pt)

		i++
		if i >= batchSize {
			if err := ztc.Ic.Write(ztc.Bp); err != nil {
				return err
			}
			i = 0
		}
	}
	if err := ztc.Ic.Write(ztc.Bp); err != nil {
		return err
	}

	return nil
}

type Ka  struct{
	Source io.Reader
	Ic     client.Client
	Bp     client.BatchPoints
}

func (ka *Ka) Parse() error {
	return nil
}
