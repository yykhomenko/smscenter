package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "smsclient"
	app.Usage = "SMPP client for the command line"
	app.Version = "0.1"
	app.Authors = []cli.Author{{"Yurii Khomenko", "yykhomenko@gmail.com"}}
	app.Flags = []cli.Flag{

		cli.StringFlag{
			Name:  "addr",
			Value: "localhost:2775",
			Usage: "Set SMPP server host:port",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Set SMPP username",
		},
		cli.StringFlag{
			Name:  "passwd",
			Usage: "Set SMPP password",
		},
		cli.BoolFlag{
			Name:  "tls",
			Usage: "Use client TLS connection",
		},
		cli.BoolFlag{
			Name:  "precaire",
			Usage: "Accept invalid TLS certificate",
		},
	}
	app.Commands = []cli.Command{
		cmdShortMessage,
		cmdQueryMessage,
	}
	app.Run(os.Args)
}
