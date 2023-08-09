package gokui

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	_ "github.com/pingcap/tidb/types/parser_driver"
)

func GenerateSelect(in string, opt GenerateSelectOptions) (string, error) {
	stmtNode, err := parse(in)
	if err != nil {
		return "", err
	}

	switch n := stmtNode.(type) {
	case *ast.CreateTableStmt:
		s := buildSelectFromCreateTable(n)
		var sql string
		if opt.NewLine {
			sql = s.multiLine()
		} else {
			sql = s.singleLine()
		}
		return sql, nil
	default:
		return "", fmt.Errorf("`%v` statement is not supported", ast.GetStmtLabel(stmtNode))
	}
}

type GenerateSelectOptions struct {
	NewLine bool
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

func (s *Select) multiLine() string {
	nl := "\n"
	ind := strings.Repeat(" ", 2)

	b := &strings.Builder{}
	b.WriteString("SELECT")
	b.WriteString(nl)

	cl := len(s.cols)
	for i, c := range s.cols {
		b.WriteString(ind)
		b.WriteString(c)
		if i < cl-1 {
			b.WriteString(",")
		}
		b.WriteString(nl)
	}

	b.WriteString("FROM")
	b.WriteString(nl)
	b.WriteString(ind)
	b.WriteString(s.table)
	b.WriteString(nl)
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
