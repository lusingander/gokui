package main

import (
	"log"
	"os"

	"github.com/lusingander/gokui/internal/gokui"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.App{
		Name: "gokui",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SQL from other SQL",
				Subcommands: []*cli.Command{
					{
						Name:   "select",
						Usage:  "Generate `SELECT`",
						Action: gokui.GenerateSelectAction,
					},
				},
			},
		},
		Before: func(ctx *cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
