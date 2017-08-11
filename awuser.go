package cicada

import (
	"os"
	"log"
	"time"
	"sort"
	"bufio"
	"regexp"
	"strings"
	"strconv"

	"github.com/urfave/cli"
)

func AwuserAction(ctx *cli.Context) error {
	start := time.Now()
	
	fi, fo, dt, cf := ctx.String("inputfile"),ctx.String("outputfile"),ctx.Int("date"),ctx.Bool("count")

	log.Println(fi, fo, dt, cf)

	if !PathExist(fi) {
		return cli.NewExitError("The file "+fi+" is not existed.", 1)
	}
	// Remove output file
	if PathExist(fo) {
		if err := os.Remove(fo); err != nil {
			return cli.NewExitError("The file "+fo+" remove failed: "+err.Error(), 1)
		}
	}

	fr, err := os.Open(fi)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		if err := fr.Close(); err != nil {
			//return cli.NewExitError(err.Error(), 1)
			panic(err)
		}
	}()
	
	fw, err := os.Create(fo)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		if err := fw.Close(); err != nil {
			//return cli.NewExitError(err.Error(), 1)
			panic(err)
		}
	}()
	w := bufio.NewWriter(fw)

	userCount := make(map[string]int, 1000000)
	scanner := bufio.NewScanner(fr)
	for scanner.Scan() {
		rawRow := scanner.Text()
		if isMatch(rawRow, dt) {
			// if count
			if cf {
				rowFields := strings.Fields(rawRow)
				lastField := strings.Split(rowFields[len(rowFields)-1], "|")
				if len(lastField) == 4 {
					userFlag := lastField[3]
					userFlag = userFlag[0:len(userFlag)-1]
					userCount[userFlag] = userCount[userFlag]+1
				}
			} else {
				rawRow = rawRow+"\n"
				if _, err := w.Write([]byte(rawRow)); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
			}
		}
	}

	if cf {
		sortUserPairs := sortUserCount(userCount)
		for _, v := range sortUserPairs {
			rowSlice := []string{v.Key, " ", strconv.Itoa(v.Value), "\n"}
			wrow := strings.Join(rowSlice, "")
			if _, err := w.Write([]byte(wrow)); err != nil {
				return cli.NewExitError(err.Error(), 1)
			}
		}
	}
	
	if err := w.Flush(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	timeTrack(start, "Total")
	
	return nil
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

func sortUserCount(userCount map[string]int) PairList {
	pl := make(PairList, len(userCount))
	i := 0
	for k, v := range userCount {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func isMatch(row string, dt int) bool {
	date := time.Now().AddDate(0,0,dt).Format("02/Jan/2006")
	reg1 := regexp.MustCompile("adview")
	reg2 := regexp.MustCompile(date)
	if reg1.MatchString(row) && reg2.MatchString(row) {
		return true
	}
	return false
}

func PathExist(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
