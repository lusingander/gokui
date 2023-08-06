package gokui

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
	"github.com/urfave/cli/v2"
)

func GenerateSelectAction(cCtx *cli.Context) error {
	var r io.Reader = os.Stdin
	var w io.Writer = os.Stdout

	bytes, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	sql := string(bytes)

	out, err := generateSelect(sql)
	if err != nil {
		return err
	}

	w.Write([]byte(out))

	return nil
}

func generateSelect(in string) (string, error) {
	stmtNode, err := parse(in)
	if err != nil {
		return "", err
	}

	switch n := stmtNode.(type) {
	case *ast.CreateTableStmt:
		sql := buildSelectFromCreateTable(n).singleLine()
		return sql, nil
	default:
		return "", fmt.Errorf("`%v` statement is not supported", ast.GetStmtLabel(stmtNode))
	}
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
