package database

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

const whereClause, fromClause, returningClause, selectClause = " where ", " from ", " returning ", "select "

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
	if where != "" {
		whereClause += where
	}
	groupByClause := ""
	if groupBy != "" {
		groupByClause = " group by " + groupBy
	}
	orderByClause := ""
	if orderBy != "" {
		orderByClause = " order by " + orderBy
	}
	if dbt.Placeholder == "@" {
		switch {
		case limit == 0 && offset == 0:
			return selectClause + columns + fromClause + tables + whereClause + groupByClause + orderByClause
		case limit > 0 && offset == 0:
			limitClause := fmt.Sprintf("top %d ", limit)
			return selectClause + limitClause + columns + fromClause + tables + whereClause + groupByClause + orderByClause
		case offset > 0:
			if orderByClause == "" {
				orderByClause = " order by (select null)"
			}
			offsetClause := fmt.Sprintf(" offset %d rows", offset)
			limitClause := fmt.Sprintf(" fetch next %d rows only", limit)
			return selectClause + columns + fromClause + tables + whereClause + groupByClause + orderByClause + offsetClause + limitClause
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
	return selectClause + columns + fromClause + tables + whereClause + groupByClause + orderByClause + limitClause + offsetClause
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
		q += returningClause + strings.Join(returning, ", ")
	}
	return q
}

func SQLUpdate(table string, columns []string, where string, dbt *DBType) string {
	whereClause := ""
	if where != "" {
		whereClause += where
	}

	stmts := lo.FilterMap(columns, func(col string, i int) (string, bool) {
		return fmt.Sprintf("%s = %s", col, dbt.PlaceholderFor(i+1)), true
	})
	return fmt.Sprintf("update %s set %s%s", table, strings.Join(stmts, ", "), whereClause)
}

func SQLUpdateReturning(table string, columns []string, where string, returned []string, dbt *DBType) string {
	q := SQLUpdate(table, columns, where, dbt)
	if len(returned) > 0 {
		q += returningClause + strings.Join(returned, ", ")
	}
	return q
}

func SQLUpsert(table string, columns []string, rows int, conflicts []string, updates []string, dbt *DBType) string {
	{{{ if .SQLServer }}}if dbt.Placeholder == "@" {
		return sqlServerUpsert(table, columns, rows, conflicts, updates, dbt)
	}
	{{{ end }}}q := SQLInsert(table, columns, rows, dbt)
	q += " on conflict (" + strings.Join(conflicts, ", ") + ") do update set "
	lo.ForEach(updates, func(x string, idx int) {
		if idx > 0 {
			q += ", "
		}
		q += fmt.Sprintf("%s = excluded.%s", x, x)
	})
	return q
}{{{ if .SQLServer }}}

func sqlServerUpsert(table string, columns []string, rows int, conflicts []string, updates []string, dbt *DBType) string {
	colNames := strings.Join(columns, ", ")
	params := make([]string, 0, rows)
	for rowIdx := range lo.Range(rows) {
		cols := make([]string, 0, len(columns))
		for colIdx := range columns {
			cols = append(cols, fmt.Sprintf("@p%d", (rowIdx*len(columns))+colIdx+1))
		}
		params = append(params, fmt.Sprintf("(%s)", strings.Join(cols, ", ")))
	}
	atSymbols := strings.Join(params, ", ")
	sourceJoin := strings.Join(lo.Map(conflicts, func(pk string, _ int) string {
		return fmt.Sprintf("%s.%s = source.%s", table, pk, pk)
	}), " and ")
	assignments := strings.Join(lo.Map(updates, func(c string, _ int) string {
		return fmt.Sprintf("%s = \"source\".%s", c, c)
	}), ", ")
	vals := strings.Join(lo.Map(updates, func(c string, _ int) string {
		return fmt.Sprintf("source.%s", c)
	}), ", ")
	sql := `merge into %s using (values %s) as source (%s) on %s when matched then update set %s when not matched then insert (%s) values (%s);`
	return fmt.Sprintf(sql, table, atSymbols, colNames, sourceJoin, assignments, colNames, vals)
}{{{ end }}}

func SQLDelete(table string, where string, _ *DBType) string {
	if strings.TrimSpace(where) == "" {
		return fmt.Sprintf("attempt to delete from [%s] with empty where clause", table)
	}
	return "delete from " + table + whereClause + where
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
