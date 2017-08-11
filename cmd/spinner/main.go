package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"
	"github.com/briandowns/spinner"
)

func main() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

	app := cli.NewApp()
	app.Name = "adview"
	app.Version = "1.0.0"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "tangshouqiang",
			Email: "tangshouqiang@tvmining.com",
		},
		cli.Author{
			Name:  "tangkefan",
			Email: "kefan@tvmining.com",
		},
	}
	app.Usage = "Adview command tools"
	app.Commands = []cli.Command{
		cli.Command{
			Name: "sync",
			//Category: "adview",
			Usage: "Synchronous adview data",
			Subcommands: cli.Commands{
				cli.Command{
					Name: "channel",
					Action: channelAction,
				},
			},
			BashComplete: func(c *cli.Context) {
				fmt.Fprintf(c.App.Writer, "--better\n")
			},
			Before: func(c *cli.Context) error {
				fmt.Fprintf(c.App.Writer, "brace for impact\n")
				s.Start()

				return nil
			},
			After: func(c *cli.Context) error {
				s.FinalMSG = "did we lose anyone?\n"
				s.Stop()
				return nil
			},
			Action: func(c *cli.Context) error {
				c.Command.FullName()
				c.Command.HasName("sync")
				fmt.Fprintf(c.App.Writer, "dodododododododododododododoo\n")
				if c.Bool("forever") {
					c.Command.Run(c)
				}
				time.Sleep(4 * time.Second)
				return nil
			},
			OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
				fmt.Fprintf(c.App.Writer, "for shame\n")
				return err
			},
		},
	}

	app.Run(os.Args)
}

func channelAction(c *cli.Context) error {
	fmt.Fprintf(c.App.Writer, ":wave: over here, eh\n")
	return nil
}
