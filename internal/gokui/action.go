package gokui

import (
	"fmt"
	"strings"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
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

func GenerateInsert(in string, opt GenerateInsertOptions) (string, error) {
	stmtNode, err := parse(in)
	if err != nil {
		return "", err
	}

	switch n := stmtNode.(type) {
	case *ast.CreateTableStmt:
		i := buildInsertFromCreateTable(n)
		var sql string
		if opt.NewLine {
			sql = i.multiLine()
		} else {
			sql = i.singleLine()
		}
		return sql, nil
	default:
		return "", fmt.Errorf("`%v` statement is not supported", ast.GetStmtLabel(stmtNode))
	}
}

type GenerateInsertOptions struct {
	NewLine bool
}

type Insert struct {
	table string
	cols  []column
}

func buildInsertFromCreateTable(stmt *ast.CreateTableStmt) *Insert {
	table := stmt.Table.Name.String()
	cols := make([]column, len(stmt.Cols))
	for i, c := range stmt.Cols {
		cols[i] = column{
			name:       c.Name.Name.String(),
			columnType: columnTypeFrom(c.Tp),
		}
	}
	return &Insert{table, cols}
}

func (i *Insert) singleLine() string {
	cns := make([]string, len(i.cols))
	dvs := make([]string, len(i.cols))
	for idx, c := range i.cols {
		cns[idx] = c.name
		dvs[idx] = c.columnType.defaultValue()
	}

	joinedCols := strings.Join(cns, ", ")
	defaultVals := strings.Join(dvs, ", ")

	b := &strings.Builder{}
	b.WriteString("INSERT INTO")
	b.WriteString(" ")
	b.WriteString(i.table)
	b.WriteString(" ")
	b.WriteString("(")
	b.WriteString(joinedCols)
	b.WriteString(")")
	b.WriteString(" ")
	b.WriteString("VALUES")
	b.WriteString(" ")
	b.WriteString("(")
	b.WriteString(defaultVals)
	b.WriteString(")")
	b.WriteString(";")
	return b.String()
}

func (i *Insert) multiLine() string {
	nl := "\n"
	ind := strings.Repeat(" ", 2)
	cl := len(i.cols)

	b := &strings.Builder{}
	b.WriteString("INSERT INTO")
	b.WriteString(" ")
	b.WriteString(i.table)
	b.WriteString(nl)

	b.WriteString("(")
	b.WriteString(nl)
	for i, c := range i.cols {
		b.WriteString(ind)
		b.WriteString(c.name)
		if i < cl-1 {
			b.WriteString(",")
		}
		b.WriteString(nl)
	}
	b.WriteString(")")
	b.WriteString(nl)

	b.WriteString("VALUES")
	b.WriteString(nl)
	b.WriteString("(")
	b.WriteString(nl)

	for i, c := range i.cols {
		b.WriteString(ind)
		b.WriteString(c.columnType.defaultValue())
		if i < cl-1 {
			b.WriteString(",")
		}
		b.WriteString(nl)
	}
	b.WriteString(")")
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

type column struct {
	name string
	columnType
}

type columnType int

const (
	unknown columnType = iota
	numberLike
	stringLike
)

func columnTypeFrom(tp *types.FieldType) columnType {
	switch tp.GetType() {
	case mysql.TypeTiny,
		mysql.TypeShort,
		mysql.TypeInt24,
		mysql.TypeLong,
		mysql.TypeLonglong,
		mysql.TypeBit,
		mysql.TypeYear,
		mysql.TypeFloat,
		mysql.TypeDouble,
		mysql.TypeNewDecimal:
		return numberLike
	default:
		return stringLike
	}
}

func (t columnType) defaultValue() string {
	switch t {
	case numberLike:
		return `0`
	case stringLike:
		return `''`
	default:
		return `''`
	}
}
