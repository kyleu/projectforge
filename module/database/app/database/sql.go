package database

import (
	"fmt"
	"strings"
)

const whereSpaces = " where "

func SQLSelect(columns string, tables string, where string, orderBy string, limit int, offset int) string {
	if columns == "" {
		columns = "*"
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = whereSpaces + where
	}

	orderByClause := ""
	if len(orderBy) > 0 {
		orderByClause = " order by " + orderBy
	}

	limitClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf(" limit %d", limit)
	}

	offsetClause := ""
	if offset > 0 {
		offsetClause = fmt.Sprintf(" offset %d", offset)
	}

	return "select " + columns + " from " + tables + whereClause + orderByClause + limitClause + offsetClause
}

func SQLSelectSimple(columns string, tables string, where ...string) string {
	return SQLSelect(columns, tables, strings.Join(where, " and "), "", 0, 0)
}

func SQLInsert(table string, columns []string, rows int, placeholder string) string {
	if rows <= 0 {
		rows = 1
	}
	colString := strings.Join(columns, ", ")
	var placeholders []string
	for i := 0; i < rows; i++ {
		var ph []string
		for idx := range columns {
			if placeholder == "?" {
				ph = append(ph, "?")
			} else {
				ph = append(ph, fmt.Sprintf("$%d", (i*len(columns))+idx+1))
			}
		}
		placeholders = append(placeholders, "("+strings.Join(ph, ", ")+")")
	}
	return fmt.Sprintf("insert into %s (%s) values %s", table, colString, strings.Join(placeholders, ", "))
}

func SQLInsertReturning(table string, columns []string, rows int, returning []string, placeholder string) string {
	q := SQLInsert(table, columns, rows, placeholder)
	if len(returning) > 0 {
		q += " returning " + strings.Join(returning, ", ")
	}
	return q
}

func SQLUpdate(table string, columns []string, where string, placeholder string) string {
	whereClause := ""
	if len(where) > 0 {
		whereClause = whereSpaces + where
	}

	stmts := make([]string, 0, len(columns))
	for i, col := range columns {
		if placeholder == "?" {
			stmts = append(stmts, fmt.Sprintf("%s = ?", col))
		} else {
			stmts = append(stmts, fmt.Sprintf("%s = $%d", col, i+1))
		}
	}
	return fmt.Sprintf("update %s set %s%s", table, strings.Join(stmts, ", "), whereClause)
}

func SQLUpdateReturning(table string, columns []string, where string, returned []string, placeholder string) string {
	q := SQLUpdate(table, columns, where, placeholder)
	if len(returned) > 0 {
		q += " returning " + strings.Join(returned, ", ")
	}
	return q
}

func SQLDelete(table string, where string) string {
	if strings.TrimSpace(where) == "" {
		return fmt.Sprintf("attempt to delete from [%s] with empty where clause", table)
	}
	return "delete from " + table + whereSpaces + where
}
