package database

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

const whereSpaces = " where "

func SQLSelect(columns string, tables string, where string, orderBy string, limit int, offset int, dbt *DBType) string {
	return SQLSelectGrouped(columns, tables, where, "", orderBy, limit, offset, dbt)
}

func SQLSelectSimple(columns string, tables string, dbt *DBType, where ...string) string {
	return SQLSelectGrouped(columns, tables, strings.Join(where, " and "), "", "", 0, 0, dbt)
}

func SQLSelectGrouped(columns string, tables string, where string, groupBy string, orderBy string, limit int, offset int, dbt *DBType) string {
	if columns == "" {
		columns = "*"
	}
	whereClause := ""
	if len(where) > 0 {
		whereClause = whereSpaces + where
	}
	groupByClause := ""
	if len(groupBy) > 0 {
		groupByClause = " group by " + groupBy
	}
	orderByClause := ""
	if len(orderBy) > 0 {
		orderByClause = " order by " + orderBy
	}
	if dbt.Placeholder == "@" {
		switch {
		case limit == 0 && offset == 0:
			return "select " + columns + " from " + tables + whereClause + groupByClause + orderByClause
		case limit > 0 && offset == 0:
			limitClause := fmt.Sprintf("top %d ", limit)
			return "select " + limitClause + columns + " from " + tables + whereClause + groupByClause + orderByClause
		case offset > 0:
			if orderByClause == "" {
				orderByClause = " order by (select null)"
			}
			offsetClause := fmt.Sprintf(" offset %d rows", offset)
			limitClause := fmt.Sprintf(" fetch next %d rows only", limit)
			return "select " + columns + " from " + tables + whereClause + groupByClause + orderByClause + offsetClause + limitClause
		}
	}
	limitClause := ""
	offsetClause := ""
	if limit > 0 {
		limitClause = fmt.Sprintf(" limit %d", limit)
	}
	if offset > 0 {
		offsetClause = fmt.Sprintf(" offset %d", offset)
	}
	return "select " + columns + " from " + tables + whereClause + groupByClause + orderByClause + limitClause + offsetClause
}

func SQLInsert(table string, columns []string, rows int, dbt *DBType) string {
	if rows <= 0 {
		rows = 1
	}
	colString := strings.Join(columns, ", ")
	var placeholders []string
	lo.Times(rows, func(i int) struct{} {
		ph := lo.Map(columns, func(_ string, idx int) string {
			return dbt.PlaceholderFor((i * len(columns)) + idx + 1)
		})
		placeholders = append(placeholders, "("+strings.Join(ph, ", ")+")")
		return struct{}{}
	})
	return fmt.Sprintf("insert into %s (%s) values %s", table, colString, strings.Join(placeholders, ", "))
}

func SQLInsertReturning(table string, columns []string, rows int, returning []string, dbt *DBType) string {
	q := SQLInsert(table, columns, rows, dbt)
	if len(returning) > 0 {
		q += " returning " + strings.Join(returning, ", ")
	}
	return q
}

func SQLUpdate(table string, columns []string, where string, dbt *DBType) string {
	whereClause := ""
	if len(where) > 0 {
		whereClause = whereSpaces + where
	}

	stmts := lo.FilterMap(columns, func(col string, i int) (string, bool) {
		return fmt.Sprintf("%s = %s", col, dbt.PlaceholderFor(i+1)), true
	})
	return fmt.Sprintf("update %s set %s%s", table, strings.Join(stmts, ", "), whereClause)
}

func SQLUpdateReturning(table string, columns []string, where string, returned []string, dbt *DBType) string {
	q := SQLUpdate(table, columns, where, dbt)
	if len(returned) > 0 {
		q += " returning " + strings.Join(returned, ", ")
	}
	return q
}

func SQLUpsert(table string, columns []string, rows int, conflicts []string, updates []string, dbt *DBType) string {
	q := SQLInsert(table, columns, rows, dbt)
	q += " on conflict (" + strings.Join(conflicts, ", ") + ") do update set "
	lo.ForEach(updates, func(x string, idx int) {
		if idx > 0 {
			q += ", "
		}
		q += fmt.Sprintf("%s = excluded.%s", x, x)
	})
	return q
}

func SQLDelete(table string, where string, _ *DBType) string {
	if strings.TrimSpace(where) == "" {
		return fmt.Sprintf("attempt to delete from [%s] with empty where clause", table)
	}
	return "delete from " + table + whereSpaces + where
}

func SQLInClause(column string, numParams int, offset int, dbt *DBType) string {
	resBuilder := strings.Builder{}
	for index := 0; index < numParams; index++ {
		if index == 0 {
			resBuilder.WriteString(column + " in (")
		} else {
			resBuilder.WriteString(", ")
		}
		resBuilder.WriteString(dbt.PlaceholderFor(index + offset + 1))
	}
	resBuilder.WriteString(")")
	return resBuilder.String()
}
