package database

import (
	"database/sql"
	"strings"
)

type SqlBuilder struct {
	query string
	args  []interface{}
}

func (sb *SqlBuilder) Append(query string, args ...interface{}) *SqlBuilder {
	sb.query += strings.TrimSpace(query) + " "
	if len(args) > 0 {
		sb.args = append(sb.args, args...)
	}
	return sb
}

func (sb *SqlBuilder) ValueList(args ...interface{}) *SqlBuilder {
	if len(args) > 0 {
		sb.query += "(" + strings.Repeat(",?", len(args))[1:] + ") "
		sb.args = append(sb.args, args...)
	} else {
		sb.query += "() "
	}
	return sb
}

func (sb *SqlBuilder) Exec() (sql.Result, error) {
	return Exec(sb.query, sb.args...)
}

func (sb *SqlBuilder) Query() (*sql.Rows, error) {
	return Query(sb.query, sb.args...)
}

func (sb *SqlBuilder) QueryRow() *sql.Row {
	return QueryRow(sb.query, sb.args...)
}
