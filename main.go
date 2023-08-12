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
		Name:  "gokui",
		Usage: "SQL utilities",
		Commands: []*cli.Command{
			{
				Name:  "generate",
				Usage: "Generate SQL from other SQL",
				Subcommands: []*cli.Command{
					{
						Name:   "select",
						Usage:  "Generate `SELECT`",
						Action: generateSelectAction,
						Flags: []cli.Flag{
							newlineFlag,
						},
					},
					{
						Name:   "insert",
						Usage:  "Generate `INSERT`",
						Action: generateInsertAction,
						Flags: []cli.Flag{
							newlineFlag,
							insertSelectFlag,
						},
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

// https://github.com/urfave/cli/issues/1122
// https://github.com/urfave/cli/pull/1128#issuecomment-1112772833
type quietBoolFlag struct {
	cli.BoolFlag
}

func (f *quietBoolFlag) String() string {
	return cli.FlagStringer(f)
}

func (f *quietBoolFlag) GetDefaultText() string {
	return ""
}

var (
	newlineFlag = &quietBoolFlag{
		cli.BoolFlag{
			Name:  "newline",
			Value: false,
			Usage: "Generate newlined output",
		},
	}
	insertSelectFlag = &quietBoolFlag{
		cli.BoolFlag{
			Name:  "insert-select",
			Value: false,
			Usage: "Generate INSERT-SELECT format output",
		},
	}
)

func generateSelectAction(cCtx *cli.Context) error {
	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout

	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sql := string(bytes)

	newLine := cCtx.Bool("newline")

	opt := gokui.GenerateSelectOptions{
		NewLine: newLine,
	}

	out, err := gokui.GenerateSelect(sql, opt)
	if err != nil {
		return err
	}

	w.Write([]byte(out))

	return nil
}

func generateInsertAction(cCtx *cli.Context) error {
	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout

	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sql := string(bytes)

	newLine := cCtx.Bool("newline")
	insertSelect := cCtx.Bool("insert-select")

	opt := gokui.GenerateInsertOptions{
		NewLine:      newLine,
		InsertSelect: insertSelect,
	}

	out, err := gokui.GenerateInsert(sql, opt)
	if err != nil {
		return err
	}

	w.Write([]byte(out))

	return nil
}
