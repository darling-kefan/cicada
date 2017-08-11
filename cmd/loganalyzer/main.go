package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/briandowns/spinner"
	"github.com/darling-kefan/cicada"
)

var spin = spinner.New(spinner.CharSets[9], 250*time.Millisecond)

var awuserCmd = cli.Command{
	Name: "awuser",
	Usage: "analysis adview user data",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "inputfile, i",
			Value: "input.log",
			Usage: "analysis file",
		},
		cli.StringFlag{
			Name:  "outputfile, o",
			Value: "output.log",
			Usage: "output write into file",
		},
		cli.IntFlag{
			Name:  "date, d",
			Usage: "The {date} log to analysis",
		},
		cli.BoolFlag{
			Name:  "count, c",
			Usage: "statistics user count",
		},
	},
	Before: func(c *cli.Context) error {
		spin.Start()
		return nil
	},
	After: func(c *cli.Context) error {
		spin.FinalMSG = "Done!\n"
		spin.Stop()
		return nil
	},
	Action: cicada.AwuserAction,
}

var awcompareCmd = cli.Command{
	Name: "awcompare",
	Usage: "compare adview user data",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "cpuprofile",
			Usage: "write cpu profile to this file",
		},
	},
	Before: func(c *cli.Context) error {
		spin.Start()
		return nil
	},
	After: func(c *cli.Context) error {
		spin.FinalMSG = "Done!\n"
		spin.Stop()
		return nil
	},
	Action: cicada.AwcompareAction,
}

func main() {
	app := cli.NewApp()
	app.Name = "loganalyzer"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "tangshouqiang",
			Email: "tangshouqiang@tvmining.com",
		},
	}
	app.Usage = "log analysis command tools"
	app.Commands = []cli.Command{
		awuserCmd,
		awcompareCmd,
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Thar be no command %q here.\n", command)
	}
	app.Run(os.Args)
}

