package cicada

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/urfave/cli"
)

func LoguserAction(ctx *cli.Context) error {
	inputFile, outputFile := ctx.String("input-file"), ctx.String("output-file")
	if inputFile == "" || outputFile == "" {
		fmt.Println("Miss flag -i or -o")
		return nil
	}
	fi, err := os.Open(inputFile)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fw, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
		if err := fw.Close(); err != nil {
			panic(err)
		}
	}()

	scanner := bufio.NewScanner(fi)
	for scanner.Scan() {
		rawRow := scanner.Text()
		if rawRow == "[]" {
			continue
		}
		rowSli := strings.Split(strings.TrimPrefix(strings.TrimSuffix(rawRow, "]"), "["), "|")
		if len(rowSli) < 4 {
			continue
		}
		fw.Write([]byte(rowSli[3]+"\n"))
	}
	if err := scanner.Err(); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	return nil
}
