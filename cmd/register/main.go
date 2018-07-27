package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"github.com/srelab/register/pkg"
	"github.com/srelab/register/pkg/g"
	"github.com/srelab/register/pkg/logger"
	"github.com/srelab/register/pkg/util"
)

func main() {
	app := &cli.App{
		Name:     g.NAME,
		Usage:    "Docker event automatically listens and synchronizes consule",
		Version:  g.VERSION,
		Compiled: time.Now(),
		Authors:  []cli.Author{{Name: g.AUTHOR, Email: g.MAIL}},
		Before: func(c *cli.Context) error {
			fmt.Fprintf(c.App.Writer, util.StripIndent(
				`
				#####  ######  ####  #  ####  ##### ###### #####
				#    # #      #    # # #        #   #      #    #
				#    # #####  #      #  ####    #   #####  #    #
				#####  #      #  ### #      #   #   #      #####
				#   #  #      #    # # #    #   #   #      #   #
				#    # ######  ####  #  ####    #   ###### #    #
			`))
			return nil
		},
		Commands: []cli.Command{
			{
				Name:  "start",
				Usage: "start a new gateway-register",
				Action: func(ctx *cli.Context) error {
					for _, flagName := range ctx.FlagNames() {
						if ctx.String(flagName) != "" {
							continue
						}

						fmt.Println(flagName + " is required")
						os.Exit(127)
					}

					g.ParseConfig(ctx)
					logger.InitLogger()

					return pkg.Start()
				},
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "concurrency", Value: 10, Usage: "concurrency number"},
					&cli.StringFlag{Name: "docker.endpoint", Value: "unix:///var/run/docker.sock", Usage: "Docker Conn EndPoint"},
					&cli.StringFlag{Name: "log.dir", Value: "./", Usage: "the log file is written to the path"},
					&cli.StringFlag{Name: "log.level", Value: "info", Usage: "valid levels: [debug, info, warn, error, fatal]"},
					&cli.StringFlag{Name: "gateway.host", Usage: "gateway server host"},
					&cli.StringFlag{Name: "gateway.port", Usage: "gateway server port"},
					&cli.StringFlag{Name: "consul.host", Usage: "consul server host"},
					&cli.StringFlag{Name: "consul.port", Usage: "consul server port"},
					&cli.StringFlag{Name: "privilege.host", Usage: "privilege server host"},
					&cli.StringFlag{Name: "privilege.port", Usage: "privilege server port"},
				},
			},
		},
	}

	app.Run(os.Args)
}
