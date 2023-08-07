package main

import (
	"io"
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
						Action: generateSelectAction,
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

func generateSelectAction(cCtx *cli.Context) error {
	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout

	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sql := string(bytes)

	out, err := gokui.GenerateSelect(sql)
	if err != nil {
		return err
	}

	w.Write([]byte(out))

	return nil
}
