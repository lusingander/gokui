package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/urfave/cli/v2"
)

func generateSelectAction(cCtx *cli.Context) error {
	return generateSelect(os.Stdin, os.Stdout)
}

func generateSelect(r io.Reader, w io.Writer) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sql := string(bytes)

	stmtNode, err := parse(sql)
	if err != nil {
		return err
	}

	switch n := stmtNode.(type) {
	case *ast.CreateTableStmt:
		sql := buildSelectFromCreateTable(n).singleLine()
		w.Write([]byte(sql))
	default:
		return fmt.Errorf("`%v` statement is not supported", ast.GetStmtLabel(stmtNode))
	}

	return nil
}

type Select struct {
	table string
	cols  []string
}

func buildSelectFromCreateTable(stmt *ast.CreateTableStmt) *Select {
	table := stmt.Table.Name.String()
	cols := make([]string, len(stmt.Cols))
	for i, c := range stmt.Cols {
		cols[i] = c.Name.Name.String()
	}
	return &Select{table, cols}
}

func (s *Select) singleLine() string {
	joinedCols := strings.Join(s.cols, ", ")

	b := &strings.Builder{}
	b.WriteString("SELECT")
	b.WriteString(" ")
	b.WriteString(joinedCols)
	b.WriteString(" ")
	b.WriteString("FROM")
	b.WriteString(" ")
	b.WriteString(s.table)
	b.WriteString(";")

	return b.String()
}

func parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNode, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmtNode, nil
}

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
