package cicada

import (
	_ "log"
	"os"
	"time"
	"bufio"
	"strings"
	"strconv"
	"runtime/pprof"

	"github.com/urfave/cli"
)

func AwcompareAction(ctx *cli.Context) error {
	cpuProfile := ctx.String("cpuprofile")
	if cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	
	if ctx.NArg() <= 1 {
		return cli.NewExitError("Error: enter at least two files", 1)
	}

	for _, f := range ctx.Args() {
		if !PathExist(f) {
			return cli.NewExitError("Error: The file "+f+" is not existed.", 1)
		}
	}

	userVisits := make(map[string][]int, 10000)
	for _, f := range ctx.Args() {
		fr, err := os.Open(f)
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
		scanner := bufio.NewScanner(fr)
		for scanner.Scan() {
			rowFields := strings.Fields(scanner.Text())
			/*if visit, err := strconv.Atoi(rowFields[0]); err != nil {
				userVisits[rowFields[1]] = append(userVisits[rowFields[1]], visit)
			}*/
			// 上面写法错误，必须先初始化，如下所示
			if _, ok := userVisits[rowFields[1]]; !ok {
				userVisits[rowFields[1]] = make([]int, 0)
			}
			if visit, err := strconv.Atoi(rowFields[0]); err == nil {
				userVisits[rowFields[1]] = append(userVisits[rowFields[1]], visit)
			}
		}

		if err := fr.Close(); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}

	// log.Println(userVisits)
	
	fw, err := os.Create("compare.log")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		if err := fw.Close(); err != nil {
			//return cli.NewExitError(err.Error(), 1)
			panic(err)
		}
	}()
	bw := bufio.NewWriter(fw)

	for k, v := range userVisits {
		rowSli := make([]string, 0)
		rowSli = append(rowSli, k)
		for _, vv := range v {
			rowSli = append(rowSli, strconv.Itoa(vv))
		}
		rowStr := strings.Join(rowSli, " ")+"\n"
		if _, err := bw.Write([]byte(rowStr)); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	}
	if err := bw.Flush(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	time.Sleep(2 * time.Second)
	
	return nil
}
