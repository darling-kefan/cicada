package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/briandowns/spinner"
	"github.com/darling-kefan/cicada"
)

var spin = spinner.New(spinner.CharSets[9], 100*time.Millisecond)
//var spin = spinner.New(spinner.CharSets[26], 500*time.Millisecond)

var spinnerCmd = cli.Command{
	Name: "spinner",
	Usage: "spinner test",
	BashComplete: func(c *cli.Context) {
		//fmt.Fprintf(c.App.Writer, "--better\n")
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
	Action: func(c *cli.Context) error {
		time.Sleep(5 * time.Second)
		return nil
	},
}

var channelCmd = cli.Command{
	Name: "channel",
	Usage: "Synchronous channel data from redis to mysql",
	Before: func(c *cli.Context) error {
		return nil
	},
	After: func(c *cli.Context) error {
		return nil
	},
	Action: cicada.ChannelAction,
}

var loguserCmd = cli.Command{
	Name: "loguser",
	Usage: "training user log data",
	Flags: []cli.Flag{
		cli.StringFlag{Name: "input-file, i"},
		cli.StringFlag{Name: "output-file, o"},
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
	Action: cicada.LoguserAction,
}


func main() {
	app := cli.NewApp()
	app.Name = "adview"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "tangshouqiang",
			Email: "tangshouqiang@tvmining.com",
		},
	}
	app.Usage = "Adview command tools"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "sync",
			Usage: "Synchronous adview data",
			Subcommands: cli.Commands{
				spinnerCmd,
				channelCmd,
			},
		},
		cli.Command{
			Name: "train",
			Usage: "train log data",
			Subcommands: cli.Commands{
				loguserCmd,
			},
		},
	}
	app.CommandNotFound = func(c *cli.Context, command string) {
		fmt.Fprintf(c.App.Writer, "Thar be no command %q here.\n", command)
	}
	app.Run(os.Args)
}
